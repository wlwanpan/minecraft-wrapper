package minecraft_1_15

import (
	"fmt"

	wrapper "github.com/wlwanpan/minecraft-wrapper"
)

type gameRuleBoolName string

func (gr gameRuleBoolName) gameRuleName() string {
	return string(gr)
}

// GameRule Names that use booleans
const (
	AnnounceAdvancements       gameRuleBoolName = "announceAdvancements"
	CommandBlockOutput         gameRuleBoolName = "commandBlockOutput"
	DisableElytraMovementCheck gameRuleBoolName = "disableElytraMovementCheck"
	DisableRaids               gameRuleBoolName = "disableRaids"
	DoDaylightCycle            gameRuleBoolName = "doDaylightCycle"
	DoEntityDrops              gameRuleBoolName = "doEntityDrops"
	DoFireTick                 gameRuleBoolName = "doFireTick"
	DoInsomnia                 gameRuleBoolName = "doInsomnia"
	DoImmediateRespawn         gameRuleBoolName = "doImmediateRespawn"
	DoLimitedCrafting          gameRuleBoolName = "doLimitedCrafting"
	DoMobLoot                  gameRuleBoolName = "doMobLoot"
	DoMobSpawning              gameRuleBoolName = "doMobSpawning"
	DoPatrolSpawning           gameRuleBoolName = "doPatrolSpawning" // new
	DoTileDrops                gameRuleBoolName = "doTileDrops"
	DoTraderSpawning           gameRuleBoolName = "doTraderSpawning" // new
	DoWeatherCycle             gameRuleBoolName = "doWeatherCycle"
	DrowningDamage             gameRuleBoolName = "drowningDamage"
	FallDamage                 gameRuleBoolName = "fallDamage"
	FireDamage                 gameRuleBoolName = "fireDamage"
	KeepInventory              gameRuleBoolName = "keepInventory"
	LogAdminCommands           gameRuleBoolName = "logAdminCommands"
	MobGriefing                gameRuleBoolName = "mobGriefing"
	NaturalRegeneration        gameRuleBoolName = "naturalRegeneration"
	ReducedDebugInfo           gameRuleBoolName = "reducedDebugInfo"
	SendCommandFeedback        gameRuleBoolName = "sendCommandFeedback"
	ShowDeathMessages          gameRuleBoolName = "showDeathMessages"
	SpectatorsGenerateChunks   gameRuleBoolName = "spectatorsGenerateChunks"
)

type gameRuleIntName string

func (gr gameRuleIntName) gameRuleName() string {
	return string(gr)
}

// GameRule Names that use integers
const (
	MaxCommandChainLength gameRuleIntName = "maxCommandChainLength"
	MaxEntityCramming     gameRuleIntName = "maxEntityCramming"
	RandomTickSpeed       gameRuleIntName = "randomTickSpeed"
	SpawnRadius           gameRuleIntName = "spawnRadius"
)

// GameRuleName accepts either boolean or integer name types
type GameRuleName interface {
	gameRuleName() string
}

// GameRule is a command used to set various rules in game
type GameRule struct {
	name GameRuleName
	bVal bool
	iVal int
}

func NewGameRuleGet(name GameRuleName) GameRule {
	return GameRule{name: name}
}

func NewGameRuleBoolean(name GameRuleName, b bool) GameRule {
	return GameRule{name: name, bVal: b}
}

func NewGameRuleInt(name GameRuleName, i int) GameRule {
	return GameRule{name: name, iVal: i}
}

// Command allows the GameRule struct to be executed as a command in game
func (c GameRule) Command() string {
	switch c.name.(type) {
	case gameRuleBoolName:
		return fmt.Sprintf("gamerule %s %t", c.name.gameRuleName(), c.bVal)
	case gameRuleIntName:
		return fmt.Sprintf("gamerule %s %d", c.name.gameRuleName(), c.iVal)
	default:
		return fmt.Sprintf("gamerule %s", c.name.gameRuleName())
	}
}

func (c GameRule) Events() []wrapper.Event {
	return []wrapper.Event{
		&GameRuleSet{},
		&GameRuleGet{},
		&wrapper.IncorrectCommandArgument{},
		&wrapper.InvalidBoolean{},
		&wrapper.InvalidInteger{},
		&wrapper.UnknownOrIncompleteCommand{},
	}
}

type GameRuleSet struct {
	name GameRuleName
	bVal *bool
	iVal *int
}

func (event *GameRuleSet) Parse(s string) bool {
	if _, err := fmt.Sscanf(s, "Gamerule %s is now set to: %T", event.bVal); err != nil {
		if _, err := fmt.Sscanf(s, "Gamerule %s is now set to: %d", event.iVal); err != nil {
			return false
		}
	}
	return true
}

type GameRuleGet struct {
	name GameRuleName
	bVal *bool
	iVal *int
}

func (event *GameRuleGet) Parse(s string) bool {
	if _, err := fmt.Sscanf(s, "Gamerule %s is currently set to: %T", event.bVal); err != nil {
		if _, err := fmt.Sscanf(s, "Gamerule %s is currently set to: %d", event.iVal); err != nil {
			return false
		}
	}
	return true
}
