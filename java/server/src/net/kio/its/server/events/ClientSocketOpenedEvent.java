package net.kio.its.server.events;

import net.kio.its.event.Cancellable;

import java.net.Socket;

public class ClientSocketOpenedEvent extends SocketEvent implements Cancellable {

    private boolean cancelled;

    public ClientSocketOpenedEvent(Socket socket) {
        super(socket);
    }

    @Override
    public boolean isCancelled() {
        return cancelled;
    }

    @Override
    public void setCancelled(boolean cancelled) {
        this.cancelled = cancelled;
    }
}
