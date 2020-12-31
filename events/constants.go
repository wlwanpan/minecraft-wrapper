package events

type EventType int

const (
	TypeNil   EventType = iota
	TypeState EventType = iota
	TypeCmd   EventType = iota
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
	Version         string = "version"
	PlayerJoined           = "player-joined"
	PlayerLeft             = "player-left"
	PlayerUUID             = "player-uuid"
	PlayerSay              = "player-say"
	PlayerDied             = "player-died"
	TimeIs                 = "time-is"
	DataGet                = "data-get"
	DataGetNoEntity        = "data-get-no-entity"
	Seed                   = "seed"
	DefaultGameMode        = "default-game-mode"
	Banned                 = "banned"
)
