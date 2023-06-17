package net.kio.its.client.events;

import net.kio.its.client.SocketWorker;
import net.kio.its.event.Event;

import java.util.UUID;

class DataEvent extends Event {

    private final SocketWorker socketWorker;
    private final String data;
    private final UUID messageUUID;

    public DataEvent(SocketWorker serverWorker, String data, UUID messageUUID) {
        this.socketWorker = serverWorker;
        this.data = data;
        this.messageUUID = messageUUID;
    }

    public SocketWorker getSocketWorker() {
        return socketWorker;
    }

    public String getData() {
        return data;
    }

    public UUID getMessageUUID() {
        return messageUUID;
    }
}
