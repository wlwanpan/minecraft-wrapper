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

func NewStateEvent(e string) StateEvent {
	return StateEvent{name: e}
}

var (
	NilEvent     StateEvent = NewStateEvent(Empty)
	StartedEvent StateEvent = NewStateEvent(Started)
	StoppedEvent StateEvent = NewStateEvent(Stopped)
	StartEvent   StateEvent = NewStateEvent(Start)
	StopEvent    StateEvent = NewStateEvent(Stop)
)

type GameEvent struct {
	id   int
	name string
	Tick int
}

func (ge GameEvent) String() string {
	return ge.name
}

func (ge GameEvent) Is(e Event) bool {
	return ge.String() == e.String()
}

func NewGameEvent(e string) GameEvent {
	gameEventCount++
	return GameEvent{
		id:   gameEventCount,
		name: e,
	}
}

var (
	TimeIsEvent GameEvent = NewGameEvent(TimeIs)
)
