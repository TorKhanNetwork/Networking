package net.kio.its.responsesystem;

import java.util.UUID;
import java.util.concurrent.TimeoutException;
import java.util.function.BiConsumer;

public class Response {

    private final IResponseWorker responseWorker;
    private final UUID msgUUID;
    private BiConsumer<? super IResponseWorker, ? super String> onReply;
    private boolean replied = false;
    private long timeout = 3000L;

    public Response(IResponseWorker ServerWorker, UUID uuid) {
        this.responseWorker = ServerWorker;
        this.msgUUID = uuid != null ? uuid : UUID.randomUUID();
    }

    public Response waitReply(boolean blockThread, BiConsumer<? super IResponseWorker, ? super String> onReply) throws TimeoutException {
        this.onReply = onReply;
        responseWorker.getResponseManager().waitResponse(this);
        if (blockThread) {
            long now = System.currentTimeMillis();
            while (!replied) {
                try {
                    Thread.sleep(100);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                if (System.currentTimeMillis() >= now + timeout) {
                    throw new TimeoutException("No response from server in " + timeout + "ms");
                }
            }
        }
        return this;
    }

    public Response setTimeout(long timeout) {
        this.timeout = timeout;
        return this;
    }

    public void acceptReply(String message) {
        this.replied = true;
        if (onReply != null) {
            onReply.accept(responseWorker, message);
        }
    }

    public UUID getMsgUUID() {
        return msgUUID;
    }
}
