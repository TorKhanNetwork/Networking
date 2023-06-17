package net.kio.its.client;

import net.kio.its.client.events.*;
import net.kio.its.event.Event;
import net.kio.its.logger.LogType;
import net.kio.its.responsesystem.IResponseWorker;
import net.kio.its.responsesystem.Response;
import net.kio.its.responsesystem.ResponseManager;
import net.kio.security.dataencryption.EncryptedRequestManager;
import net.kio.security.dataencryption.KeysGenerator;

import java.io.*;
import java.net.Socket;
import java.net.SocketException;
import java.security.KeyFactory;
import java.security.NoSuchAlgorithmException;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;
import java.util.UUID;

public class SocketWorker extends Thread implements IResponseWorker {

    private final Client client;

    private final String targetIp;
    private final int targetPort;
    private Socket socket;
    private boolean isConnected;
    private boolean isAuthenticated;

    private DataOutputStream dataOutputStream;
    private DataInputStream dataInputStream;

    private final KeysGenerator keysGenerator;
    private final ResponseManager responseManager;

    private String commandPrefix;
    private String requestSeparator;

    private boolean connectionProtocol;

    public SocketWorker(Client client, String targetIp, int targetPort) {
        this.client = client;
        this.targetIp = targetIp;
        this.targetPort = targetPort;
        this.keysGenerator = new KeysGenerator();
        this.responseManager = new ResponseManager();
        setName("Socket Worker #" + getId());
    }

    public void startWorker() {
        try {
            connectSocket();
            start();
        } catch (IOException exception) {
            exception.printStackTrace();
        }
    }

    public void stopWorker() {
        if (isConnected) {
            sendCommand("disconnect");
            interrupt();
        }
    }

    public void connectSocket() throws IOException {
        if (socket != null && isConnected) throw new SocketException("Socket is already connected");
        else {
            socket = new Socket(targetIp, targetPort);
            socket.setSoTimeout(0);
            dataOutputStream = new DataOutputStream(socket.getOutputStream());
            dataInputStream = new DataInputStream(socket.getInputStream());
            isConnected = true;
        }
    }

    public void disconnectSocket() throws IOException {
        if (socket == null) throw new SocketException("Socket isn't connected");
        else {
            stopWorker();
            socket.close();
            isConnected = false;
            interrupt();
            socket = null;
        }
    }

    @Override
    public void run() {
        handleServerSocket();
    }

    private void handleServerSocket() {
        try {
            BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(dataInputStream));
            startConnectionProtocol();
            String line;
            while (!Thread.interrupted() && (line = bufferedReader.readLine()) != null) {
                onLineRead(line);
            }
            disconnectSocket();
        } catch (SocketException e) {
            client.getLogger().log(LogType.CRITICAL, getName() + " - Unlegit disconnection from server " + targetIp + " on port " + targetPort);
            ServerSocketClosedEvent event = new ServerSocketClosedEvent(socket, this);
            client.getEventsManager().callEvent(event);
            interrupt();
        } catch (IOException exception) {
            exception.printStackTrace();
        }
    }

    private void onLineRead(String line) {

        String decryptedData;
        if (keysGenerator.getSecretKey() != null && (decryptedData = EncryptedRequestManager.decrypt(line, keysGenerator)) != null) {
            onDataReceived(new String[]{decryptedData, line}, true);
        } else {
            onDataReceived(new String[]{line}, false);
        }


    }

    private void onDataReceived(String[] data, boolean encrypted) {

        UUID msgUUID;

        if (data[0].startsWith("response:")) {
            msgUUID = UUID.fromString(data[0].substring(9, 45));
            data[0] = data[0].substring(45);
            client.getLogger().log(LogType.DEBUG, getName() + " - Response received for UUID " + msgUUID);
            responseManager.onResponseReceived(msgUUID, data[0]);
        } else {
            msgUUID = UUID.fromString(data[0].substring(0, 36));
            data[0] = data[0].substring(36);
        }

        if (commandPrefix != null && data[0].startsWith(commandPrefix)) {
            onCommandReceived(data[0], msgUUID, encrypted);
        } else {
            client.getLogger().log(LogType.DEBUG, getName() + " - Data received (encrypted=" + encrypted + ", uuid=" + msgUUID + ") : " + data[0]);

            Event event = new RawDataReceivedEvent(this, data[0], msgUUID);
            if (encrypted) {
                event = new EncryptedDataReceivedEvent(this, data[1], data[0], msgUUID);
            }
            client.getEventsManager().callEvent(event);


            if (connectionProtocol) {
                if (!onConnectionProtocolDataReceived(data[0])) {
                    client.getLogger().log(LogType.CRITICAL, "Unable to connect to the server");
                    stopWorker();
                }
            }


        }

    }

    private void onCommandReceived(String command, UUID msgUUID, boolean encrypted) {
        client.getLogger().log(LogType.DEBUG, getName() + " -  Command Received (encrypted=" + encrypted + ", uuid=" + msgUUID + ") : " + command);
        CommandReceivedEvent event = new CommandReceivedEvent(this, command, commandPrefix, requestSeparator, msgUUID);
        client.getEventsManager().callEvent(event);
    }


    private void startConnectionProtocol() {
        connectionProtocol = true;
        sendData("version:" + Client.getVersion(), false);
    }

    private boolean onConnectionProtocolDataReceived(String data) {

        if (data.startsWith("version:")) {
            boolean version = Boolean.parseBoolean(data.substring(8));
            if (!version) return false;
            else {
                sendData("macAddress:" + Client.getMacAddress(), false);
            }
        } else if (data.startsWith("publicKey:")) {
            String array = data.replace("publicKey:", "");
            try {
                keysGenerator.setPublicKey(KeyFactory.getInstance("RSA").generatePublic(new X509EncodedKeySpec(Base64.getDecoder().decode(array))));
                keysGenerator.generateKeys(true, false);
                sendData("secretKey:" + EncryptedRequestManager.encryptSecretKey(keysGenerator), false);
            } catch (InvalidKeySpecException | NoSuchAlgorithmException e) {
                e.printStackTrace();
                return false;
            }

        } else if (data.startsWith("separator:")) {
            requestSeparator = data.replace("separator:", "");
        } else if (data.startsWith("commandPrefix:")) {
            commandPrefix = data.replace("commandPrefix:", "");
        } else if (data.startsWith("connected:")) {
            boolean connected = Boolean.parseBoolean(data.replace("connected:", ""));
            isAuthenticated = connected;
            connectionProtocol = !connected;
            if (connected) {
                ConnectionProtocolSuccessEvent event = new ConnectionProtocolSuccessEvent(this);
                client.getEventsManager().callEvent(event);
            }
            return connected;
        } else {
            return false;
        }

        return true;
    }

    public Response sendData(String data) {
        return sendData(data, (UUID) null, true);
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
        client.getLogger().log(LogType.DEBUG, getName() + " - Data sent (encrypted=" + encrypt + ", uuid=" + response.getMsgUUID() + ") : " + rawData);
        return response;
    }

    public Response sendCommand(String command, String... args) {
        return sendData(commandPrefix + command + requestSeparator + String.join(requestSeparator, args));
    }

    public boolean isConnected() {
        return isConnected;
    }

    public boolean isAuthenticated() {
        return isAuthenticated;
    }

    public String getTargetIp() {
        return targetIp;
    }

    public int getTargetPort() {
        return targetPort;
    }

    public ResponseManager getResponseManager() {
        return responseManager;
    }
}
