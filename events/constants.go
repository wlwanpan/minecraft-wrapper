package events

type EventType int

const (
	TypeNil   EventType = iota
	TypeState EventType = iota
	TypeGame  EventType = iota
)

const (
	Empty string = "empty"
)

// State related events that has a direct effect on the wrapper state.
const (
	Started  = "started"
	Stopped  = "stopped"
	Starting = "starting"
	Stopping = "stopping"
	Saving   = "saving"
	Saved    = "saved"
)

// Game related events that provide player/server related information.
const (
	Version      string = "version"
	PlayerJoined string = "player-joined"
	PlayerLeft          = "player-left"
	PlayerUUID          = "player-uuid"
	PlayerSay           = "player-say"
	TimeIs              = "time-is"
	DataGet             = "data-get"
)
