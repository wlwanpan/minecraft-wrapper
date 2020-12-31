package wrapper

import (
	"testing"
	"time"

	"github.com/wlwanpan/minecraft-wrapper/events"
)

func TestBasicEventsQueue(t *testing.T) {
	eqm := newEventsQueue()
	c1 := eqm.get("event-queue-1")
	c2 := eqm.get("event-queue-2")

	go func() {
		time.Sleep(10 * time.Millisecond)
		eqm.push(events.NewGameEvent("event-queue-1"))
	}()

	select {
	case <-c1:
		// Event received from queue 1, all good!
	case ev2 := <-c2:
		t.Errorf("received event %s from the wrong queue", ev2.String())
	case <-time.After(50 * time.Millisecond):
		t.Errorf("timeout: no events received")
	}
}
