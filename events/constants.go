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
	Stopped  string = "stopped"
	Starting string = "starting"
	Stopping string = "stopping"
	Saving   string = "saving"
	Saved    string = "saved"
)

// Game related events that provide player/server related information.
const (
	Banned           string = "banned"
	BanList          string = "ban-list"
	BanListEntry     string = "ban-list-entry"
	DataGet          string = "data-get"
	DataGetNoEntity  string = "data-get-no-entity"
	DefaultGameMode  string = "default-game-mode"
	Difficulty       string = "difficulty"
	ExperienceAdd    string = "experience-add"
	ExperienceQuery  string = "experience-query"
	Give             string = "give"
	NoPlayerFound    string = "no-player-found"
	PlayerJoined     string = "player-joined"
	PlayerLeft       string = "player-left"
	PlayerUUID       string = "player-uuid"
	PlayerSay        string = "player-say"
	PlayerDied       string = "player-died"
	Kicked           string = "kicked"
	Seed             string = "seed"
	ServerOverloaded string = "server-overloaded"
	TimeIs           string = "time-is"
	UnknownItem      string = "unknown-item"
	Version          string = "version"
	WhisperTo        string = "whisper-to"
)
