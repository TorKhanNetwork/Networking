package net.kio.its.client.events;

import net.kio.its.client.SocketWorker;

import java.net.Socket;

public class ServerDisconnectEvent extends SocketEvent {

    private final SocketWorker socketWorker;

    public ServerDisconnectEvent(Socket socket, SocketWorker socketWorker) {
        super(socket);
        this.socketWorker = socketWorker;
    }

    public SocketWorker getSocketWorker() {
        return socketWorker;
    }

}
