package net.kio.its.event;

import java.util.ArrayList;
import java.util.List;

public class ThreadedEventCaller extends Thread {

    private final EventsManager eventsManager;
    private final List<Event> eventsToCall;

    public ThreadedEventCaller(EventsManager eventsManager) {
        this.eventsManager = eventsManager;
        this.eventsToCall = new ArrayList<>();
    }

    public void callEvent(Event event) {
        synchronized (eventsToCall) {
            eventsToCall.add(event);
            synchronized (this) {
                notify();
            }
        }
    }

    @Override
    public void run() {
        while (!Thread.interrupted()) {
            try {
                if (eventsToCall.size() == 0) {
                    synchronized (this) {
                        wait();
                    }
                }
                synchronized (eventsToCall) {
                    List<Event> eventsCalled = new ArrayList<>();
                    for (Event event : eventsToCall) {
                        eventsManager.callEvent0(event);
                        eventsCalled.add(event);
                    }
                    eventsToCall.removeAll(eventsCalled);
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }
    }
}
