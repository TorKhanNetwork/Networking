package events_system

import "reflect"

type RegisteredListener struct {
	listener      Listener
	eventExecutor EventExecutor
}

func NewRegisteredListener(listener Listener, eventExecutor EventExecutor) RegisteredListener {
	return RegisteredListener{
		listener:      listener,
		eventExecutor: eventExecutor,
	}
}

func (registeredListener *RegisteredListener) CallEvent(event *Event) {
	if f := reflect.ValueOf(event).FieldByName("Cancel"); f == reflect.ValueOf(nil) || !f.Bool() {
		registeredListener.eventExecutor(registeredListener.listener, event)
	}
}
