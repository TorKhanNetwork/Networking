package events_system

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

func (registeredListener *RegisteredListener) CallEvent(event Event) {
	var i interface{} = event
	if c, ok := i.(Cancellable); !ok || !c.IsCancelled() {
		registeredListener.eventExecutor(registeredListener.listener, event)
	}
}
