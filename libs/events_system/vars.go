package events_system

type Event interface {
}

type Listener interface {
}

type EventExecutor func(Listener, Event) *Event

type EventPriority int

const (
	LOWEST  EventPriority = 0
	LOW     EventPriority = 1
	NORMAL  EventPriority = 2
	HIGH    EventPriority = 3
	HIGHEST EventPriority = 4
	MONITOR EventPriority = 5
)
