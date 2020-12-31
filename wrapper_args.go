package wrapper

type GameMode string

const (
	Survival  GameMode = "survival"
	Creative  GameMode = "creative"
	Adventure GameMode = "adventure"
	Spectator GameMode = "spectator"
)

type GameDifficulty string

const (
	Easy     GameDifficulty = "easy"
	Hard     GameDifficulty = "hard"
	Normal   GameDifficulty = "normal"
	Peaceful GameDifficulty = "peaceful"
)
