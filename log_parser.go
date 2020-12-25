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
	events.Starting: regexp.MustCompile(`Starting minecraft server version (.*)`),
	events.Stopping: regexp.MustCompile(`Stopping (.*) server`),
	events.Saving:   regexp.MustCompile(`Saving the game`),
	events.Saved:    regexp.MustCompile(`Saved (?s)(.*)`),
}

var gameEventToRegex = map[string]*regexp.Regexp{
	events.PlayerJoined: regexp.MustCompile(`(?s)(.*) joined the game`),
	events.PlayerLeft:   regexp.MustCompile(`(?s)(.*) left the game`),
	events.PlayerUUID:   regexp.MustCompile(`UUID of player (?s)(.*) is (?s)(.*)`),
	events.PlayerSay:    regexp.MustCompile(`<(?s)(.*)> (?s)(.*)`),
	events.TimeIs:       regexp.MustCompile(`The time is (?s)(.*)`),
	events.DataGet:      regexp.MustCompile(`(?s)(.*) has the following (entity|block|storage) data: (.*)`),
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
		case events.TimeIs:
			return handleTimeEvent(matches)
		case events.PlayerUUID:
			return handlePlayerUUIDEvent(matches, tick)
		case events.PlayerSay:
			return handlePlayerSayEvent(matches, tick)
		case events.DataGet:
			return handleDataGet(matches, tick)
		default:
			gameEvent := events.NewGameEvent(e)
			gameEvent.Tick = tick
			return gameEvent, events.TypeGame
		}
	}
	return events.NilEvent, events.TypeNil
}

func handleTimeEvent(matches []string) (events.GameEvent, events.EventType) {
	tickStr := matches[1]
	tick, _ := strconv.Atoi(tickStr)
	timeEvent := events.NewGameEvent(events.TimeIs)
	timeEvent.Tick = tick
	return timeEvent, events.TypeGame
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

func handleDataGet(matches []string, tick int) (events.GameEvent, events.EventType) {
	dgEvent := events.NewGameEvent(events.DataGet)
	dgEvent.Tick = tick
	dgEvent.Data = map[string]string{
		"player_name": matches[1],
		"data_type":   matches[2],
		"data_raw":    matches[3], // TODO: Need to unmarshall str -> struct | interface{}
	}
	return dgEvent, events.TypeGame
}
