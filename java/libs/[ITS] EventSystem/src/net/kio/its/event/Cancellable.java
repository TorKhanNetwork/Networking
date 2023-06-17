package net.kio.its.event;

public interface Cancellable {
    boolean isCancelled();

    void setCancelled(boolean cancelled);
}
