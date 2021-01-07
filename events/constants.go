package events

type EventType int

const (
	TypeNil EventType = iota
	TypeState
	TypeCmd
	TypeGame
)

const (
	Empty string = "empty"
)

// State related events that has a direct effect on the wrapper state.
const (
	Started  string = "started"
	Stopped         = "stopped"
	Starting        = "starting"
	Stopping        = "stopping"
	Saving          = "saving"
	Saved           = "saved"
)

// Game related events that provide player/server related information.
const (
	Banned           string = "banned"
	BanList                 = "ban-list"
	BanListEntry            = "ban-list-entry"
	DataGet                 = "data-get"
	DataGetNoEntity         = "data-get-no-entity"
	DefaultGameMode         = "default-game-mode"
	Difficulty              = "difficulty"
	ExperienceAdd           = "experience-add"
	ExperienceQuery         = "experience-query"
	Give                    = "give"
	NoPlayerFound           = "no-player-found"
	PlayerJoined            = "player-joined"
	PlayerLeft              = "player-left"
	PlayerUUID              = "player-uuid"
	PlayerSay               = "player-say"
	PlayerDied              = "player-died"
	Kicked                  = "kicked"
	Seed                    = "seed"
	ServerOverloaded        = "server-overloaded"
	TimeIs                  = "time-is"
	UnknownItem             = "unknown-item"
	Version                 = "version"
	WhisperTo               = "whisper-to"
)
