package net.kio.its.server;


import net.kio.its.event.Event;
import net.kio.its.logger.LogType;
import net.kio.its.responsesystem.IResponseWorker;
import net.kio.its.responsesystem.Response;
import net.kio.its.responsesystem.ResponseManager;
import net.kio.its.server.events.*;
import net.kio.security.dataencryption.EncryptedRequestManager;
import net.kio.security.dataencryption.KeysGenerator;

import java.io.*;
import java.net.Socket;
import java.net.SocketException;
import java.util.Arrays;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

public class ServerWorker extends Thread implements IResponseWorker {

    private final Server server;
    private final Socket socket;
    private boolean isConnected;
    private boolean connectionProtocol;

    private final DataOutputStream dataOutputStream;
    private final BufferedReader bufferedReader;

    private final KeysGenerator keysGenerator;
    private final ResponseManager responseManager;
    private final String requestSeparator;
    private String macAddress;
    private final String commandPrefix;

    public ServerWorker(Server server, Socket socket, KeysGenerator keysGenerator) throws IOException {
        KeysGenerator keysGenerator1;
        this.server = server;
        this.socket = socket;
        socket.setSoTimeout(0);
        this.dataOutputStream = new DataOutputStream(socket.getOutputStream());
        this.bufferedReader = new BufferedReader(new InputStreamReader(socket.getInputStream()));
        try {
            keysGenerator1 = keysGenerator.clone();
        } catch (CloneNotSupportedException e) {
            keysGenerator1 = new KeysGenerator();
        }
        this.keysGenerator = keysGenerator1;
        this.responseManager = new ResponseManager();
        this.requestSeparator = Server.generateRequestSeparator();
        this.commandPrefix = "!";
        setName("Server Worker #" + getId());
    }

    public void stopWorker() {
        if (isConnected) {
            sendData("disconnect", (!connectionProtocol));
            interrupt();
        }
    }

    public void disconnectSocket() throws IOException {
        if (socket == null) throw new SocketException("Socket isn't connected");
        else {
            stopWorker();
            socket.close();
            isConnected = false;
            interrupt();
        }
    }

    @Override
    public void run() {
        handleClientSocket();
    }

    public void handleClientSocket() {
        try {
            startConnectionProtocol();
            String line;
            while (!Thread.interrupted() && (line = bufferedReader.readLine()) != null) {
                onLineRead(line);
            }
            disconnectSocket();
        } catch (SocketException e) {
            server.getLogger().log(LogType.CRITICAL, getName() + " - Unlegit disconnection from client " + macAddress);
            ClientSocketClosedEvent event = new ClientSocketClosedEvent(socket, this);
            server.getEventsManager().callEvent(event);
        } catch (IOException exception) {
            exception.printStackTrace();
        }
    }

    private void onLineRead(String line) {
        String decryptedData = null;
        if (keysGenerator.getSecretKey() != null && (decryptedData = EncryptedRequestManager.decrypt(line, keysGenerator)) != null) {
            onDataReceived(new String[]{decryptedData, line}, true);
        } else {
            if (keysGenerator.getSecretKey() != null)
                System.out.println(Arrays.toString(keysGenerator.getSecretKey().getEncoded()));
            onDataReceived(new String[]{line}, false);
        }
    }

    private void onDataReceived(String[] data, boolean encrypted) {

        UUID msgUUID;

        if (data[0].startsWith("response:")) {
            msgUUID = UUID.fromString(data[0].substring(9, 45));
            data[0] = data[0].substring(45);
            server.getLogger().log(LogType.DEBUG, getName() + " - Response received for UUID " + msgUUID);
            responseManager.onResponseReceived(msgUUID, data[0]);
        } else {
            msgUUID = UUID.fromString(data[0].substring(0, 36));
            data[0] = data[0].substring(36);
        }

        if (data[0].startsWith(commandPrefix)) {
            onCommandReceived(data[0], encrypted, msgUUID);
        } else {
            server.getLogger().log(LogType.DEBUG, getName() + " - Data received (encrypted=" + encrypted + ", uuid=" + msgUUID + ") : " + data[0]);
            Event event = new RawDataReceivedEvent(this, data[0], msgUUID);
            if (encrypted) {
                event = new EncryptedDataReceivedEvent(this, data[1], data[0], msgUUID);
            }
            server.getEventsManager().callEvent(event);

            if (connectionProtocol) {
                if (!onConnectionProtocolDataReceived(data[0])) {
                    server.getLogger().log(LogType.CRITICAL, "Unable to establish a connection with the client");
                    stopWorker();
                }
            }
        }

    }

    private void onCommandReceived(String command, boolean encrypted, UUID msgUUID) {
        server.getLogger().log(LogType.DEBUG, getName() + " -  Command Received (encrypted=" + encrypted + ", uuid=" + msgUUID + ") : " + command);
        CommandReceivedEvent event = new CommandReceivedEvent(this, command, commandPrefix, requestSeparator, msgUUID);
        if (event.getCommand().equalsIgnoreCase("disconnect")) {
            server.getEventsManager().callEvent(new ClientDisconnectEvent(this));
            Thread.currentThread().interrupt();
        } else {
            server.getEventsManager().callEvent(event);
        }
    }

    private void startConnectionProtocol() {
        connectionProtocol = true;
    }

    private boolean onConnectionProtocolDataReceived(String data) {
        if (data.startsWith("version:")) {
            try {
                List<Integer> clientVersion = Arrays.stream(data.replace("version:", "").replace(".", "!").split("!")).map(Integer::parseInt).collect(Collectors.toList());
                List<Integer> serverVersion = Arrays.stream(Server.getVersion().replace(".", "!").split("!")).map(Integer::parseInt).collect(Collectors.toList());
                boolean b = clientVersion.get(0).equals(serverVersion.get(0)) && clientVersion.get(1).equals(serverVersion.get(1));
                sendData("version:" + b, false);
                if (!b) {
                    return false;
                }
            } catch (NumberFormatException ignored) {
            }
            sendData("commandPrefix:" + commandPrefix, false);
            sendData("separator:" + requestSeparator, false);
        } else if (data.startsWith("macAddress:")) {
            macAddress = data.replace("macAddress:", "").toUpperCase();
            if (keysGenerator.getPublicKey() == null) keysGenerator.generateKeys(false, true);
            sendData("publicKey:" + keysGenerator.getStringPublicKey(), false);
        } else if (data.startsWith("secretKey:")) {
            EncryptedRequestManager.decryptSecretKey(data.replace("secretKey:", ""), keysGenerator);
            boolean connected = keysGenerator.getSecretKey() != null;
            sendData("connected:" + connected, false);
            connectionProtocol = !connected;
            if (connected) {
                ConnectionProtocolSuccessEvent event = new ConnectionProtocolSuccessEvent(this);
                server.getEventsManager().callEvent(event);
            }
            return connected;
        } else {
            return false;
        }
        return true;
    }


    public Response sendData(String data) {
        return sendData(data, true);
    }

    public Response sendData(String data, boolean encrypt) {
        return sendData(data, (UUID) null, encrypt);
    }

    public Response sendData(String data, String responseUUID) {
        return sendData(data, UUID.fromString(responseUUID), true);
    }

    public Response sendData(String data, String responseUUID, boolean encrypt) {
        return sendData(data, UUID.fromString(responseUUID), encrypt);
    }

    public Response sendData(String data, UUID responseUUID, boolean encrypt) {
        String rawData = data;
        Response response = new Response(this, responseUUID);
        OutputStreamWriter outputStreamWriter = new OutputStreamWriter(dataOutputStream);
        PrintWriter printWriter = new PrintWriter(outputStreamWriter, true);
        data = (responseUUID != null ? "response:" : "") + response.getMsgUUID() + data;
        if (encrypt) {
            data = EncryptedRequestManager.encrypt(data, keysGenerator);
        }
        printWriter.println(data);
        server.getLogger().log(LogType.DEBUG, getName() + " - Data sent (encrypted=" + encrypt + ", uuid=" + response.getMsgUUID() + ") : " + rawData);
        return response;
    }

    public Response sendCommand(String command, String... args) {
        return sendData(commandPrefix + command + requestSeparator + String.join(requestSeparator, args));
    }

    public String getClientMacAddress() {
        return macAddress;
    }

    public Socket getClientSocket() {
        return socket;
    }

    public boolean isConnected() {
        return isConnected;
    }

    public ResponseManager getResponseManager() {
        return responseManager;
    }
}
