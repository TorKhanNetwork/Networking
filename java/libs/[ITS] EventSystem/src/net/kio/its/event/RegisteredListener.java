package net.kio.its.event;

public class RegisteredListener {

    private final Listener listener;
    private final EventPriority eventPriority;
    private final EventExecutor executor;

    public RegisteredListener(Listener listener, EventPriority eventPriority, EventExecutor eventExecutor) {
        this.listener = listener;
        this.eventPriority = eventPriority;
        this.executor = eventExecutor;
    }

    public EventPriority getEventPriority() {
        return eventPriority;
    }

    public void callEvent(Event event) {
        if (!(event instanceof Cancellable) || !(((Cancellable) event).isCancelled())) {
            executor.execute(listener, event);
        }
    }
}
