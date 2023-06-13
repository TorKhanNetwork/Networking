package events_system

type Event interface {
}

type Listener interface {
}

type Cancellable interface {
	IsCancelled() bool
	SetCancelled(bool)
}

type EventExecutor func(Listener, Event)

type EventPriority int

const (
	LOWEST  EventPriority = 0
	LOW     EventPriority = 1
	NORMAL  EventPriority = 2
	HIGH    EventPriority = 3
	HIGHEST EventPriority = 4
	MONITOR EventPriority = 5
)
