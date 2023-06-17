package net.kio.its.event;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.Comparator;
import java.util.HashMap;
import java.util.LinkedHashSet;
import java.util.Map;
import java.util.stream.Collectors;

public class EventsManager {

    private final Map<Class<? extends Event>, LinkedHashSet<RegisteredListener>> registeredListeners;
    private final ThreadedEventCaller threadedEventCaller;

    public EventsManager() {
        this.registeredListeners = new HashMap<>();
        this.threadedEventCaller = new ThreadedEventCaller(this);
    }

    public void registerListener(Listener listener) {
        if (!threadedEventCaller.isAlive()) threadedEventCaller.start();
        Method[] methods = listener.getClass().getDeclaredMethods();
        for (Method method : methods) {
            if (method.getAnnotation(EventHandler.class) != null && method.getParameters().length == 1 && Event.class.isAssignableFrom(method.getParameterTypes()[0])) {
                final Class<? extends Event> eventClass = method.getParameterTypes()[0].asSubclass(Event.class);
                method.setAccessible(true);
                LinkedHashSet<RegisteredListener> eventSet = registeredListeners.get(eventClass);
                if (eventSet == null) {
                    eventSet = new LinkedHashSet<>();
                }
                EventExecutor executor = (eventExecutorListener, event) -> {
                    if (eventClass.isAssignableFrom(event.getClass())) {
                        try {
                            method.invoke(listener, event);
                        } catch (IllegalAccessException | InvocationTargetException e) {
                            e.printStackTrace();
                        }
                    }
                };
                RegisteredListener registeredListener = new RegisteredListener(listener, method.getAnnotation(EventHandler.class).priority(), executor);
                eventSet.add(registeredListener);
                eventSet = eventSet.stream().sorted(Comparator.comparing(RegisteredListener::getEventPriority)).collect(Collectors.toCollection(LinkedHashSet::new));
                registeredListeners.put(eventClass, eventSet);
            }
        }

    }

    void callEvent0(Event event) {
        if (registeredListeners.containsKey(event.getClass())) {
            registeredListeners.get(event.getClass()).forEach(registeredListener -> {
                registeredListener.callEvent(event);
            });
        }
    }

    public void callEvent(Event event) {
        if (!(event instanceof Cancellable)) {
            threadedEventCaller.callEvent(event);
        } else {
            callEvent0(event);
        }
    }


}
