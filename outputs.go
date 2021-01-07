package wrapper

// DataGetOutput represents the structured data logged from the
// '/data get entity' command. Some fields might not be of the
// right or precise type since the decoder will coerse any value
// to either a string, int or float64 for simplicity.
type DataGetOutput struct {
	Brain                  Brain
	HurtByTimestamp        int
	SleepTimer             int
	SpawnForced            int
	Attributes             []Attribute
	Invulnerable           int
	FallFlying             int
	PortalCooldown         int
	AbsorptionAmount       float64
	Abilities              Abilities
	FallDistance           float64
	RecipeBook             RecipeBook
	DeathTime              int
	XpSeed                 int
	XpTotal                int
	UUID                   []interface{} // Technically []int, cc issue in decoder.go#Lexer.buildStr()
	PlayerGameType         int
	SeenCredits            int
	Motion                 []float64
	Health                 float64
	FoodSaturationLevel    float64
	Air                    int
	OnGround               int
	Dimension              string
	Rotation               []float64
	XpLevel                int
	Score                  int
	Pos                    []float64
	PreviousPlayerGameType int
	Fire                   int
	XpP                    float64
	EnderItems             []interface{}
	DataVersion            int
	FoodLevel              int
	FoodExhaustionLevel    float64
	HurtTime               int // TODO: support native time.Time?
	SelectedItemSlot       int
	Inventory              Inventory
	FoodTickTimer          int
}

type Brain map[string]map[string]interface{}

type Attribute struct {
	Name string
	Base float64
}

type Abilities map[string]float64

type Inventory []map[string]interface{}

type RecipeBook struct {
	Recipes                             []string
	ToBeDisplayed                       []string
	IsBlastingFurnaceFilteringCraftable int
	IsSmokerGuiOpen                     int
	IsFilteringCraftable                int
	IsFurnaceGuiOpen                    int
	IsGuiOpen                           int
	IsFurnaceFilteringCraftable         int
	IsBlastingFurnaceGuiOpen            int
	IsSmokerFilteringCraftable          int
}

type Player struct {
	Name string
	UUID string
}
