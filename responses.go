package wrapper

type strToStrMap map[string]string

type strToFloat map[string]float64

type DataGetResponse struct {
	Brain            map[string]map[string]interface{} `json:"Brain"`
	HurtByTimestamp  int                               `json:"HurtByTimestamp"`
	SleepTimer       int                               `json:"SleepTimer"`
	SpawnForced      int                               `json:"SpawnForced"`
	Attributes       []strToStrMap                     `json:"Attributes"`
	Invulnerable     int                               `json:"Invulnerable"`
	FallFlying       int                               `json:"FallFlying"`
	PortalCooldown   int                               `json:"PortalCooldown"`
	AbsorptionAmount float64                           `json:"AbsorptionAmount"`
	Abilities        strToFloat                        `json:"abilities"`
	FallDistance     float64                           `json:"FallDistance"`
	Pos              []float64                         `json:"Pos"`
}
