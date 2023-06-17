package net.kio.its.server.events;

import net.kio.its.server.ServerWorker;

import java.net.Socket;

public class ClientSocketClosedEvent extends SocketEvent {

    private final ServerWorker serverWorker;

    public ClientSocketClosedEvent(Socket socket, ServerWorker serverWorker) {
        super(socket);
        this.serverWorker = serverWorker;
    }

    public ServerWorker getServerWorker() {
        return serverWorker;
    }
}
