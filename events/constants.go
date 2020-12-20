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
	Started string = "started"
	Stopped        = "stopped"
	Start          = "start"
	Stop           = "stop"
)

const (
	PlayerJoined string = "player-joined"
	PlayerLeft          = "player-left"
	TimeIs              = "time-is"
	Saved               = "saved"
)
