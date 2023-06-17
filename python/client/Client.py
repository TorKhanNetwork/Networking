import platform
import random
from time import sleep

from getmac import get_mac_address as gma

from libs.log_system import Logger
from SocketWorker import SocketWorker


class Client:
    def __init__(self, name: str, debug: bool = False):
        self.name = name
        self.version = "1.0.0"
        self.pythonVersion = "python-" + str(platform.python_version())
        self.macAddress = self.getMacAddress()
        self.debug = debug
        self.socketWorkerList = []
        self.logger = Logger.Logger(self)

    def getMacAddress(self):
        macAddress = gma()
        if macAddress is None:
            return self.randomMacAddress()
        return macAddress

    def randomMacAddress(self):
        macAddr = []
        for _ in range(6):
            randStr = "".join(random.sample("0123456789abcdef", 2))
            macAddr.append(randStr)
        return ":".join(macAddr)

    def addSocketWorker(self, ip: str, port: int):
        socketWorkerObj = SocketWorker(self, ip, port)
        self.socketWorkerList.append(socketWorkerObj)
        return socketWorkerObj

    def getVersion(self):
        return self.version

    def getPythonVersion(self):
        return self.pythonVersion


client = Client("Main Client", True)
socketWorker = client.addSocketWorker("localhost", 40000)
socketWorker.startWorker()
while (not socketWorker.authenticated):
    sleep(1)
inp = None
while (inp != "exit"):
    if inp is not None:
        socketWorker.sendCommand(inp)
    inp = input()
socketWorker.stopWorker()
socketWorker.join()
