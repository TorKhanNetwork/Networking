package net.kio.its.server.events;

import net.kio.its.event.Event;
import net.kio.its.server.ServerWorker;

public class ConnectionProtocolSuccessEvent extends Event {

    private final ServerWorker serverWorker;

    public ConnectionProtocolSuccessEvent(ServerWorker serverWorker) {
        this.serverWorker = serverWorker;
    }

    public ServerWorker getServerWorker() {
        return serverWorker;
    }
}
