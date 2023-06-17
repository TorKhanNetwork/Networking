package net.kio.its.client.events;

import net.kio.its.client.SocketWorker;

import java.util.UUID;

public class RawDataReceivedEvent extends DataEvent {
    public RawDataReceivedEvent(SocketWorker socketWorker, String data, UUID messageUUID) {
        super(socketWorker, data, messageUUID);
    }
}
