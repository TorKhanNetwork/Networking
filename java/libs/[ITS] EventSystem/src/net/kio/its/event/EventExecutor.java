package net.kio.its.event;

public interface EventExecutor {

    void execute(Listener listener, Event event);

}
