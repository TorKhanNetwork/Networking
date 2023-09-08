package events_system

import (
	"reflect"
	"strings"
)

type EventsManager struct {
	registeredListeners map[interface{}][]RegisteredListener
}

func NewEventsManager() EventsManager {
	e := EventsManager{
		registeredListeners: make(map[interface{}][]RegisteredListener),
	}
	return e
}

func (eventsManager *EventsManager) RegisterListener(listener Listener) {
	lType := reflect.TypeOf(listener)
	for i := 0; i < lType.NumMethod(); i++ {
		method := lType.Method(i)
		if method.Type.NumIn() == 2 && strings.HasSuffix(method.Type.In(1).String(), "Event") {
			eType := method.Type.In(1).Name()
			executor := func(l Listener, e Event) *Event {
				method.Func.Call([]reflect.Value{reflect.ValueOf(l), reflect.ValueOf(e)})
				return &e
			}
			registeredListener := NewRegisteredListener(listener, executor)
			eventsManager.registeredListeners[eType] = append(eventsManager.registeredListeners[eType], registeredListener)
		}
	}
}

func (eventsManager *EventsManager) CallEvent(event *Event) {
	if f := reflect.ValueOf(*event).FieldByName("Cancel"); f != reflect.ValueOf(nil) {
		eventsManager.CallEvent0(event)
	} else {
		go eventsManager.CallEvent0(event)
	}
}

func (eventsManager *EventsManager) CallEvent0(event *Event) {
	if listeners := eventsManager.registeredListeners[reflect.TypeOf(*event).Name()]; listeners != nil {
		for _, l := range listeners {
			l.CallEvent(event)
		}
	}
}
