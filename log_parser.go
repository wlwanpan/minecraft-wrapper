package wrapper

import (
	"log"
	"regexp"

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

type LogParser func(string) (events.Event, int)

var stateEventToRegexp = map[string]*regexp.Regexp{
	events.Started: regexp.MustCompile(`Done (?s)(.*)! For help, type "help"\n`),
	events.Start:   regexp.MustCompile(`Starting minecraft server version (.*)\n`),
	events.Stop:    regexp.MustCompile(`Stopping (.*) server\n`),
}

var gameEventToRegex = map[string]*regexp.Regexp{
	events.PlayerJoined: regexp.MustCompile(`(?s)(.*) joined the game\n`),
	events.PlayerLeft:   regexp.MustCompile(`(?s)(.*) left the game\n`),
	events.TimeIs:       regexp.MustCompile(`The time is (?s)(.*)\n`),
	events.Saved:        regexp.MustCompile(`Saved (?s)(.*)\n`),
}

func LogParserFunc(line string) (events.Event, int) {
	ll := ParseToLogLine(line)
	log.Println(ll.output)

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
			return events.NewGameEvent(e, 0), events.TypeGame
		}
	}
	return events.NilEvent, events.TypeNil
}
