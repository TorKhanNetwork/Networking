package net.kio.its.event;

public abstract class Event {

    private String name;

    public String getName() {
        return name != null ? name : this.getClass().getSimpleName();
    }
}
