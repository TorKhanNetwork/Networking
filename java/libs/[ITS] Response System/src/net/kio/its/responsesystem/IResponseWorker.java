package net.kio.its.responsesystem;

import java.util.UUID;

public interface IResponseWorker {

    ResponseManager getResponseManager();

    Response sendData(String data);

    Response sendData(String data, boolean encrypt);

    Response sendData(String data, String responseUUID);

    Response sendData(String data, String responseUUID, boolean encrypt);

    Response sendData(String data, UUID responseUUID, boolean encrypt);

    Response sendCommand(String command, String... args);
}
