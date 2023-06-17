package net.kio.its.server.events;

import net.kio.its.server.ServerWorker;

import java.util.UUID;

public class EncryptedDataReceivedEvent extends DataEvent {

    private final String decryptedData;

    public EncryptedDataReceivedEvent(ServerWorker serverWorker, String data, String decryptedData, UUID messageUUID) {
        super(serverWorker, data, messageUUID);
        this.decryptedData = decryptedData;
    }

    public String getDecryptedData() {
        return decryptedData;
    }
}
