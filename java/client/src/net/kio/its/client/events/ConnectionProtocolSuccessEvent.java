package net.kio.its.client.events;

import net.kio.its.client.SocketWorker;
import net.kio.its.event.Event;

public class ConnectionProtocolSuccessEvent extends Event {

    private final SocketWorker socketWorker;

    public ConnectionProtocolSuccessEvent(SocketWorker serverWorker) {
        this.socketWorker = serverWorker;
    }

    public SocketWorker getSocketWorker() {
        return socketWorker;
    }
}
