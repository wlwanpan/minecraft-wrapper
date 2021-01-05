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
	NilEvent      = NewStateEvent(Empty)
	StartedEvent  = NewStateEvent(Started)
	StoppedEvent  = NewStateEvent(Stopped)
	StartingEvent = NewStateEvent(Starting)
	StoppingEvent = NewStateEvent(Stopping)
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
	NilGameEvent       = NewGameEvent(Empty)
	VersionEvent       = NewGameEvent(Version)
	TimeIsEvent        = NewGameEvent(TimeIs)
	DataGetEvent       = NewGameEvent(DataGet)
	NoPlayerFoundEvent = NewGameEvent(NoPlayerFound)
	UnknownItemEvent   = NewGameEvent(UnknownItem)
	PlayerLeftEvent    = NewGameEvent(PlayerLeft)
	PlayerUUIDEvent    = NewGameEvent(PlayerUUID)
)
