package events_system

import (
	"reflect"
	"strings"
)

type EventsManager struct {
	registeredListeners map[interface{}][]RegisteredListener
	threadedEventCaller ThreadedEventCaller
}

func NewEventsManager() EventsManager {
	t := ThreadedEventCaller{}
	e := EventsManager{
		registeredListeners: make(map[interface{}][]RegisteredListener),
		threadedEventCaller: t,
	}
	t.Start(e)
	return e
}

func (eventsManager *EventsManager) RegisterListener(listener Listener) {
	lType := reflect.TypeOf(listener)
	for i := 0; i < lType.NumMethod(); i++ {
		method := lType.Method(i)
		if method.Type.NumIn() == 2 && strings.HasSuffix(method.Type.In(1).Name(), "Event") {
			eType := method.Type.In(1)
			executor := func(l Listener, e Event) {
				method.Func.Call([]reflect.Value{reflect.ValueOf(l), reflect.ValueOf(e)})
			}
			registeredListener := NewRegisteredListener(listener, executor)
			eventsManager.registeredListeners[eType.Name()] = append(eventsManager.registeredListeners[eType.Name()], registeredListener)
		}
	}
}

func (eventsManager *EventsManager) CallEvent(event *Event) {
	if f := reflect.ValueOf(*event).FieldByName("Cancel"); f != reflect.ValueOf(nil) {
		eventsManager.threadedEventCaller.CallEvent(event)
	} else {
		eventsManager.CallEvent0(event)
	}
}

func (eventsManager *EventsManager) CallEvent0(event *Event) {
	if listeners := eventsManager.registeredListeners[reflect.TypeOf(event).Name()]; listeners != nil {
		for _, l := range listeners {
			l.CallEvent(event)
		}
	}
}
