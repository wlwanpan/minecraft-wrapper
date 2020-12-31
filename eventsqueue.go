package wrapper

import (
	"sync"

	"github.com/wlwanpan/minecraft-wrapper/events"
)

type eventsQueue struct {
	mu sync.RWMutex
	q  map[string]chan events.GameEvent
}

func newEventsQueue() *eventsQueue {
	return &eventsQueue{
		q: make(map[string]chan events.GameEvent),
	}
}

func (eq *eventsQueue) get(e string) <-chan events.GameEvent {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	_, ok := eq.q[e]
	if !ok {
		eq.q[e] = make(chan events.GameEvent)
	}
	return eq.q[e]
}

func (eq *eventsQueue) push(ev events.GameEvent) {
	eq.mu.RLock()
	defer eq.mu.RUnlock()

	c, ok := eq.q[ev.String()]
	if !ok {
		// No channel is registered, means no cmd awaits a response from this events queue.
		// we can hence discard/ignore the event...
		return
	}
	select {
	case c <- ev:
	default:
	}
}
