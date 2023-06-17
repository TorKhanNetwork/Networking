package net.kio.its.responsesystem;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

public class ResponseManager {


    private final Map<UUID, Response> responsesWaitingList;

    public ResponseManager() {
        this.responsesWaitingList = new HashMap<>();
    }

    public void waitResponse(Response response) {
        responsesWaitingList.put(response.getMsgUUID(), response);
    }

    public void onResponseReceived(UUID uuid, String message) {
        if (responsesWaitingList.containsKey(uuid)) {
            responsesWaitingList.get(uuid).acceptReply(message);
            responsesWaitingList.remove(uuid);
        }
    }
}
