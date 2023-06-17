package net.kio.its.client.events;

import net.kio.its.client.SocketWorker;

import java.net.Socket;

public class ServerSocketClosedEvent extends SocketEvent {

    private final SocketWorker socketWorker;

    public ServerSocketClosedEvent(Socket socket, SocketWorker socketWorker) {
        super(socket);
        this.socketWorker = socketWorker;
    }

    public SocketWorker getSocketWorker() {
        return socketWorker;
    }
}
