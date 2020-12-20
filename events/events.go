package events

var (
	gameEventCount int = 0
)

type Event interface {
	String() string
	Is(Event) bool
}

type StateEvent struct {
	name string
}

func (se StateEvent) String() string {
	return se.name
}

func (se StateEvent) Is(ev Event) bool {
	return se.String() == ev.String()
}

func NewStateEvent(e string) Event {
	return StateEvent{name: e}
}

var (
	NilEvent     Event = NewStateEvent(Empty)
	StartedEvent Event = NewStateEvent(Started)
	StoppedEvent Event = NewStateEvent(Stopped)
	StartEvent   Event = NewStateEvent(Start)
	StopEvent    Event = NewStateEvent(Stop)
)

type GameEvent struct {
	id   int
	name string
	tick int
}

func (ge GameEvent) String() string {
	return ge.name
}

func (ge GameEvent) Is(e Event) bool {
	return ge.String() == e.String()
}

func NewGameEvent(e string, tick int) GameEvent {
	gameEventCount++
	return GameEvent{
		id:   gameEventCount,
		name: e,
		tick: tick,
	}
}
