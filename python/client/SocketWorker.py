import socket as socketlib
import threading
import uuid

from libs.response_system import IResponseWorker, Response
from SocketStreamReader import SocketStreamReader
from libs.data_encryption import KeyGenerator, EncryptionManager


class SocketWorker(IResponseWorker.IResponseWorker):
    def __init__(self, client, ip, port):
        IResponseWorker.IResponseWorker.__init__(self)
        self.client = client
        self.ip = ip
        self.port = port
        self.socket = None
        self.connected = False
        self.authenticated = False
        self.reader = None
        self.keyGenerator = KeyGenerator.KeyGenerator()
        self.commandPrefix = None
        self.requestSeparator = None
        self.connectionProtocol = False

    def startWorker(self):
        self.__connectSocket()
        self.start()
        self.name = "Socket Worker #" + str(self.native_id)

    def stopWorker(self):
        if self.connected:
            self.sendData("disconnect", encrypt=False)
            self.interrupt()

    def __connectSocket(self):
        self.socket = socketlib.socket()
        self.socket.connect((self.ip, self.port))
        self.reader = SocketStreamReader(self.socket)
        self.connected = True
        self.client.logger.log("Socket is connected")
        # socket.send(bytes("env:" + self.client.getPythonVersion() + "\n", 'utf-8'))
        # socket.send(b"version:1.1.0\n")
        # socket.close()

    def __disconnectSocket(self):
        if self.socket is not None:
            self.socket.close()
            self.connected = False
            self.interrupt()
            self.socket = None

    def run(self) -> None:
        self.__handleServerSocket()

    def __handleServerSocket(self):
        self.__startConnectionProtocol()
        t = threading.currentThread()
        while getattr(t, "do_run", True):
            line = self.reader.readLine().replace("\n", "").replace("\r", "")
            if line is not None:
                self.__onLineRead(line)
            else:
                break
        self.__disconnectSocket()

    def __onLineRead(self, line: str):
        decryptedData = None
        if self.keyGenerator.secretKey is not None:
            decryptedData = EncryptionManager.decrypt(
                line, self.keyGenerator)
        if decryptedData is not None:
            self.__onDataReceived([decryptedData, line], True)
        else:
            self.__onDataReceived([line], False)

    def __onDataReceived(self, data: list, encrypted: bool):
        if data[0].startswith("response:"):
            msgUUID = uuid.UUID(data[0][9:45])
            data[0] = data[0][45:]
            self.client.logger.log(
                self.name + " - Response received for UUID " + str(msgUUID))
            self.responseManager.onResponseReceived(msgUUID, data[0])
        else:
            msgUUID = uuid.UUID(data[0][:36])
            data[0] = data[0][36:]

        if self.commandPrefix is not None and str(data[0]).startswith(self.commandPrefix):
            self.__onCommandReceived(data[0], msgUUID, encrypted)
        else:
            self.client.logger.log(self.name + " - Data Received (encrypted=" + str(encrypted) + ", uuid=" + str(msgUUID) + ") : "
                                   + data[0])

            # todo event system

            if self.connectionProtocol and not (self.__onConnectionProtocolDataReceived(data[0])):
                self.client.logger.log(
                    self.name + " - Unable to connect to the server")
                self.stopWorker()

    def __onCommandReceived(self, command: str, msgUUID: uuid.UUID, encrypted: bool):
        self.client.logger.log(self.name + " - Command Received (encrypted=" + str(
            encrypted) + ", uuid=" + str(msgUUID) + ") : "
            + command)

        # todo event system

    def __startConnectionProtocol(self):
        self.connectionProtocol = True
        self.sendData("version:" + self.client.getVersion(), encrypt=False)

    def __onConnectionProtocolDataReceived(self, data: str) -> bool:
        if data.startswith("version:"):
            version = data[8:] == "true"
            if not version:
                return False
            self.sendData("macAddress:" +
                          self.client.getMacAddress(), encrypt=False)
        elif data.startswith("publicKey:"):
            self.keyGenerator.publicKey = EncryptionManager.decryptPublicKey(
                data[10:])
            self.keyGenerator.generateKeys(True, False)
            self.sendData(
                "secretKey:" + EncryptionManager.encryptSecretKey(self.keyGenerator), encrypt=False)
        elif data.startswith("separator:"):
            self.requestSeparator = data[10:]
        elif data.startswith("commandPrefix:"):
            self.commandPrefix = data[14:]
        elif data.startswith("connected:"):
            connected = data[10:] == "true"
            self.authenticated = connected
            self.connectionProtocol = not connected
            if connected:
                # todo event system
                pass
            return connected
        else:
            return False
        return True

    def sendData(self, data: str, msgUUID: uuid.UUID = None, encrypt: bool = True):
        rawData = str(data)
        response = Response.Response(self, msgUUID)
        if msgUUID is not None:
            data = "response:" + data
        data = str(response.msgUUID) + data
        if encrypt:
            data = EncryptionManager.encrypt(
                data, self.keyGenerator)
        self.socket.send(bytes(data + "\n", 'utf-8'))
        self.client.logger.log(self.name + " - Data sent (encrypted=" + str(encrypt) + ", uuid=" + str(response.msgUUID) + ") : "
                               + rawData)
        return response

    def sendCommand(self, command: str, args: tuple = None):
        return self.sendData(self.commandPrefix + command + self.requestSeparator + (self.requestSeparator.join(args) if args is not None else ""))

    def interrupt(self):
        setattr(threading.currentThread(), 'do_run', False)
