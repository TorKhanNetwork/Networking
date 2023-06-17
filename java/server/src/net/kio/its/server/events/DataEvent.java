package net.kio.its.server.events;

import net.kio.its.event.Event;
import net.kio.its.server.ServerWorker;

import java.util.UUID;

class DataEvent extends Event {

    private final ServerWorker serverWorker;
    private final String data;
    private final UUID messageUUID;

    public DataEvent(ServerWorker serverWorker, String data, UUID messageUUID) {
        this.serverWorker = serverWorker;
        this.data = data;
        this.messageUUID = messageUUID;
    }

    public ServerWorker getServerWorker() {
        return serverWorker;
    }

    public String getData() {
        return data;
    }

    public UUID getMessageUUID() {
        return messageUUID;
    }
}
