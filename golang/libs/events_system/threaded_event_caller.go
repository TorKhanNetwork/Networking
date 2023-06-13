package events_system

type ThreadedEventCaller struct {
	eventsToCall []Event
	exit         chan int
	wait         chan int
	lock         bool
}

func NewThreadedEventCaller() ThreadedEventCaller {
	return ThreadedEventCaller{
		eventsToCall: make([]Event, 0),
		exit:         make(chan int),
		wait:         make(chan int),
		lock:         false,
	}
}

func (threadedEventCaller *ThreadedEventCaller) Start(eventsManager EventsManager) {
	go func() {
		for {
			select {
			case <-threadedEventCaller.exit:
				return
			default:
				if len(threadedEventCaller.eventsToCall) == 0 {
					<-threadedEventCaller.wait
				}
				threadedEventCaller.Sync()
				list := make([]Event, len(threadedEventCaller.eventsToCall))
				threadedEventCaller.lock = true
				copy(threadedEventCaller.eventsToCall, list)
				threadedEventCaller.lock = false
				for _, e := range list {
					eventsManager.CallEvent0(e)
				}
				threadedEventCaller.Sync()
				threadedEventCaller.lock = true
				threadedEventCaller.eventsToCall = threadedEventCaller.eventsToCall[len(threadedEventCaller.eventsToCall):]
				threadedEventCaller.lock = false
			}
		}
	}()
}

func (threadedEventCaller *ThreadedEventCaller) CallEvent(event Event) {
	go func() {
		threadedEventCaller.Sync()
		threadedEventCaller.lock = true
		threadedEventCaller.eventsToCall = append(threadedEventCaller.eventsToCall, event)
		threadedEventCaller.lock = false
		if len(threadedEventCaller.eventsToCall) == 1 {
			threadedEventCaller.wait <- 1
		}
	}()
}

func (threadedEventCaller *ThreadedEventCaller) Sync() {
	for {
		if !threadedEventCaller.lock {
			return
		}
	}
}
