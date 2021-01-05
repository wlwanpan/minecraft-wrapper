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
	NilEvent      StateEvent = NewStateEvent(Empty)
	StartedEvent  StateEvent = NewStateEvent(Started)
	StoppedEvent  StateEvent = NewStateEvent(Stopped)
	StartingEvent StateEvent = NewStateEvent(Starting)
	StoppingEvent StateEvent = NewStateEvent(Stopping)
)

type GameEvent struct {
	id   int
	Name string
	Tick int
	Data map[string]string
}

func (ge GameEvent) String() string {
	return ge.Name
}

func (ge GameEvent) Is(e Event) bool {
	return ge.String() == e.String()
}

func NewGameEvent(e string) GameEvent {
	gameEventCount++
	return GameEvent{
		id:   gameEventCount,
		Name: e,
	}
}

var (
	NilGameEvent       GameEvent = NewGameEvent(Empty)
	VersionEvent       GameEvent = NewGameEvent(Version)
	TimeIsEvent        GameEvent = NewGameEvent(TimeIs)
	DataGetEvent       GameEvent = NewGameEvent(DataGet)
	NoPlayerFoundEvent GameEvent = NewGameEvent(NoPlayerFound)
	UnknownItemEvent   GameEvent = NewGameEvent(UnknownItem)
)
