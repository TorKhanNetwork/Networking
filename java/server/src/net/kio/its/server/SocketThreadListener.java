package net.kio.its.server;

import net.kio.its.logger.LogType;

import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;

public class SocketThreadListener extends Thread {

    private final Server server;
    private final ServerSocket serverSocket;

    public SocketThreadListener(Server server) throws IOException {
        this.server = server;
        this.serverSocket = new ServerSocket(server.getPort());
    }

    @Override
    public void run() {

        try {
            server.getLogger().log("SocketThreadListener#run() -> Listening for clients connections");
            while (!Thread.interrupted()) {
                Socket socket = serverSocket.accept();
                server.getLogger().log(LogType.DEBUG, "New socket accepted from " + socket.getInetAddress().getHostAddress() + " on port " + socket.getPort());
                server.handleSocketConnection(socket);
            }
            serverSocket.close();
        } catch (Exception exception) {
            exception.printStackTrace();
        }

    }

    public ServerSocket getServerSocket() {
        return serverSocket;
    }
}
