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

const (
	Started  string = "started"
	Stopped         = "stopped"
	Starting        = "starting"
	Stopping        = "stopping"
	Saving          = "saving"
	Saved           = "saved"
)

const (
	PlayerJoined string = "player-joined"
	PlayerLeft          = "player-left"
	PlayerUUID          = "player-uuid"
	PlayerSay           = "player-say"
	TimeIs              = "time-is"
	DataGet             = "data-get"
)
