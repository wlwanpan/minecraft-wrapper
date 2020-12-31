package wrapper

import (
	"regexp"
	"strconv"

	"github.com/wlwanpan/minecraft-wrapper/events"
)

var logRegex = regexp.MustCompile(`(\[[0-9:]*\]) \[([A-z(-| )#0-9]*)\/([A-z #]*)\]: (.*)`)

type LogLine struct {
	timestamp  string
	threadName string
	level      string
	output     string
}

func ParseToLogLine(line string) *LogLine {
	matches := logRegex.FindAllStringSubmatch(line, 4)
	return &LogLine{
		timestamp:  matches[0][1],
		threadName: matches[0][2],
		level:      matches[0][3],
		output:     matches[0][4],
	}
}

type LogParser func(string, int) (events.Event, events.EventType)

var stateEventToRegexp = map[string]*regexp.Regexp{
	events.Started:  regexp.MustCompile(`Done (?s)(.*)! For help`),
	events.Starting: regexp.MustCompile(`Starting Minecraft server on (.*)`),
	events.Stopping: regexp.MustCompile(`Stopping (.*) server`),
	events.Saving:   regexp.MustCompile(`Saving the game`),
	events.Saved:    regexp.MustCompile(`Saved (?s)(.*)`),
}

var gameEventToRegex = map[string]*regexp.Regexp{
	events.PlayerJoined: regexp.MustCompile(`(?s)(.*) joined the game`),
	events.PlayerLeft:   regexp.MustCompile(`(?s)(.*) left the game`),
	// TODO: There is an insane amount of death messages: https://minecraft.gamepedia.com/Death_messages, support all?
	events.PlayerDied:      regexp.MustCompile(`(?s)(.*) (was shot|was pummeled|drowned|blew up|was blown up|was killed by|hit the ground|fell|was slain|suffocated)(.*)`),
	events.PlayerUUID:      regexp.MustCompile(`UUID of player (?s)(.*) is (?s)(.*)`),
	events.PlayerSay:       regexp.MustCompile(`<(?s)(.*)> (?s)(.*)`),
	events.Version:         regexp.MustCompile(`Starting minecraft server version (.*)`),
	events.TimeIs:          regexp.MustCompile(`The time is (?s)(.*)`),
	events.DataGet:         regexp.MustCompile(`(?s)(.*) has the following (entity|block|storage) data: (.*)`),
	events.DataGetNoEntity: regexp.MustCompile(`No (entity|block|storage) was found`),
	events.Seed:            regexp.MustCompile(`Seed: (.*)`),
}

func LogParserFunc(line string, tick int) (events.Event, events.EventType) {
	ll := ParseToLogLine(line)
	if ll.output == "" {
		return events.NilEvent, events.TypeNil
	}

	for e, reg := range stateEventToRegexp {
		if reg.MatchString(ll.output) {
			return events.NewStateEvent(e), events.TypeState
		}
	}
	for e, reg := range gameEventToRegex {
		matches := reg.FindStringSubmatch(ll.output)
		if matches == nil {
			continue
		}
		switch e {
		case events.PlayerJoined:
			return handlePlayerJoined(matches, tick)
		case events.PlayerLeft:
			return handlePlayerLeft(matches, tick)
		case events.PlayerDied:
			return handlePlayerDied(matches, tick)
		case events.PlayerUUID:
			return handlePlayerUUIDEvent(matches, tick)
		case events.PlayerSay:
			return handlePlayerSayEvent(matches, tick)
		case events.Version:
			return handleVersionEvent(matches)
		case events.TimeIs:
			return handleTimeEvent(matches)
		case events.DataGet:
			return handleDataGet(matches)
		case events.DataGetNoEntity:
			return handleDataGetNoEntity(matches)
		case events.Seed:
			return handleSeed(matches)
		default:
			gameEvent := events.NewGameEvent(e)
			gameEvent.Tick = tick
			return gameEvent, events.TypeGame
		}
	}
	return events.NilEvent, events.TypeNil
}

func handlePlayerJoined(matches []string, tick int) (events.GameEvent, events.EventType) {
	pjEvent := events.NewGameEvent(events.PlayerJoined)
	pjEvent.Tick = tick
	pjEvent.Data = map[string]string{
		"player_name": matches[1],
	}
	return pjEvent, events.TypeGame
}

func handlePlayerLeft(matches []string, tick int) (events.GameEvent, events.EventType) {
	plEvent := events.NewGameEvent(events.PlayerLeft)
	plEvent.Tick = tick
	plEvent.Data = map[string]string{
		"player_name": matches[1],
	}
	return plEvent, events.TypeGame
}

func handlePlayerDied(matches []string, tick int) (events.GameEvent, events.EventType) {
	pdEvent := events.NewGameEvent(events.PlayerDied)
	pdEvent.Tick = tick
	pdEvent.Data = map[string]string{
		"player_name":   matches[1],
		"death_by":      matches[2],
		"death_details": "",
	}
	if len(matches) >= 4 {
		pdEvent.Data["death_details"] = matches[3]
	}
	return pdEvent, events.TypeGame
}

func handlePlayerUUIDEvent(matches []string, tick int) (events.GameEvent, events.EventType) {
	puEvent := events.NewGameEvent(events.PlayerUUID)
	puEvent.Tick = tick
	puEvent.Data = map[string]string{
		"player_name": matches[1],
		"player_uuid": matches[2],
	}
	return puEvent, events.TypeGame
}

func handlePlayerSayEvent(matches []string, tick int) (events.GameEvent, events.EventType) {
	psEvent := events.NewGameEvent(events.PlayerSay)
	psEvent.Tick = tick
	psEvent.Data = map[string]string{
		"player_name":    matches[1],
		"player_message": matches[2],
	}
	return psEvent, events.TypeGame
}

func handleVersionEvent(matches []string) (events.GameEvent, events.EventType) {
	versionEvent := events.NewGameEvent(events.Version)
	versionEvent.Data = map[string]string{
		"version": matches[1],
	}
	return versionEvent, events.TypeCmd
}

func handleTimeEvent(matches []string) (events.GameEvent, events.EventType) {
	tickStr := matches[1]
	tick, _ := strconv.Atoi(tickStr)
	timeEvent := events.NewGameEvent(events.TimeIs)
	timeEvent.Tick = tick
	return timeEvent, events.TypeCmd
}

func handleDataGet(matches []string) (events.GameEvent, events.EventType) {
	dgEvent := events.NewGameEvent(events.DataGet)
	dgEvent.Data = map[string]string{
		"player_name": matches[1],
		"data_type":   matches[2],
		"data_raw":    matches[3],
	}
	return dgEvent, events.TypeCmd
}

func handleDataGetNoEntity(matches []string) (events.GameEvent, events.EventType) {
	dgEvent := events.NewGameEvent(events.DataGet)
	dgEvent.Data = map[string]string{
		"error_message": matches[0],
	}
	return dgEvent, events.TypeCmd
}

func handleSeed(matches []string) (events.GameEvent, events.EventType) {
	sdEvent := events.NewGameEvent(events.Seed)
	sdEvent.Data = map[string]string{
		"data_raw": matches[1],
	}
	return sdEvent, events.TypeCmd
}
