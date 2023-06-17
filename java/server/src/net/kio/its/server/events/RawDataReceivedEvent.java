package net.kio.its.server.events;

import net.kio.its.server.ServerWorker;

import java.util.UUID;

public class RawDataReceivedEvent extends DataEvent {
    public RawDataReceivedEvent(ServerWorker serverWorker, String data, UUID messageUUID) {
        super(serverWorker, data, messageUUID);
    }
}
