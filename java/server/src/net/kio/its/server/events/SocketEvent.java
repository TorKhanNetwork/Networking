package net.kio.its.server.events;

import net.kio.its.event.Event;

import java.net.Socket;

public class SocketEvent extends Event {

    private final Socket socket;

    public SocketEvent(Socket socket) {
        this.socket = socket;
    }

    public Socket getSocket() {
        return socket;
    }

}
