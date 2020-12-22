package wrapper

import (
	"log"
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

type LogParser func(string, int) (events.Event, int)

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
	events.TimeIs:       regexp.MustCompile(`The time is (?s)(.*)`),
}

func LogParserFunc(line string, t int) (events.Event, int) {
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
		if reg.MatchString(ll.output) {
			gameEvent := events.NewGameEvent(e)
			gameEvent.Tick = t
			if e == events.TimeIs {
				tickStr := reg.FindStringSubmatch(ll.output)[1]
				tick, err := strconv.Atoi(tickStr)
				if err != nil {
					log.Println("error parsing game tick: ", tickStr)
				} else {
					gameEvent.Tick = tick
				}
			}
			return gameEvent, events.TypeGame
		}
	}
	return events.NilEvent, events.TypeNil
}
