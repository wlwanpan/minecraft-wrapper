package wrapper

// DataGetOutput represents the structured data logged from the
// '/data get entity' command. Some fields might not be of the
// right or precise type since the decoder will coerse any value
// to either a string, int or float64 for simplicity.
type DataGetOutput struct {
	Brain                  Brain         `json:"Brain"`
	HurtByTimestamp        int           `json:"HurtByTimestamp"`
	SleepTimer             int           `json:"SleepTimer"`
	SpawnForced            int           `json:"SpawnForced"`
	Attributes             []interface{} `json:"Attributes"`
	Invulnerable           int           `json:"Invulnerable"`
	FallFlying             int           `json:"FallFlying"`
	PortalCooldown         int           `json:"PortalCooldown"`
	AbsorptionAmount       float64       `json:"AbsorptionAmount"`
	Abilities              Abilities     `json:"abilities"`
	FallDistance           float64       `json:"FallDistance"`
	RecipeBook             RecipeBook    `json:"recipeBook"`
	DeathTime              int           `json:"DeathTime"`
	XpSeed                 int           `json:"XpSeed"`
	XpTotal                int           `json:"XpTotal"`
	UUID                   []interface{} `json:"UUID"` // Technically []int, cc issue in decoder.go#Lexer.buildStr()
	PlayerGameType         int           `json:"playerGameType"`
	SeenCredits            int           `json:"seenCredits"`
	Motion                 []float64     `json:"Motion"`
	Health                 float64       `json:"Health"`
	FoodSaturationLevel    float64       `json:"foodSaturationLevel"`
	Air                    int           `json:"Air"`
	OnGround               int           `json:"OnGround"`
	Dimension              string        `json:"Dimension"`
	Rotation               []float64     `json:"Rotation"`
	XpLevel                int           `json:"XpLevel"`
	Score                  int           `json:"Score"`
	Pos                    []float64     `json:"Pos"`
	PreviousPlayerGameType int           `json:"previousPlayerGameType"`
	Fire                   int           `json:"Fire"`
	XpP                    float64       `json:"XpP"`
	EnderItems             []interface{} `json:"EnderItems"`
	DataVersion            int           `json:"DataVersion"`
	FoodLevel              int           `json:"foodLevel"`
	FoodExhaustionLevel    float64       `json:"foodExhaustionLevel"`
	HurtTime               int           `json:"HurtTime"` // TODO: support native time.Time?
	SelectedItemSlot       int           `json:"SelectedItemSlot"`
	Inventory              Inventory     `json:"Inventory"`
	FoodTickTimer          int           `json:"foodTickTimer"`
}

type Brain map[string]map[string]interface{}

type Abilities map[string]float64

type Inventory []map[string]interface{}

type RecipeBook map[string]interface{}

type Player struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}
