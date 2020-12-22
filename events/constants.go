package events

const (
	TypeNil   = iota
	TypeState = iota
	TypeGame  = iota
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
	TimeIs              = "time-is"
)
