package wrapper

import (
	"log"
	"regexp"
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

type LogParser func(string) Event

var eventToRegexp = map[Event]*regexp.Regexp{
	StartedEvent: regexp.MustCompile(`Done (?s)(.*)! For help, type "help"`),
	StartEvent:   regexp.MustCompile(`Starting minecraft server version (.*)`),
	StopEvent:    regexp.MustCompile(`Stopping (.*) server`),
}

func LogParserFunc(line string) Event {
	ll := ParseToLogLine(line)
	log.Println(ll.output)

	if ll.output == "" {
		return EmptyEvent
	}

	for event, reg := range eventToRegexp {
		if reg.MatchString(ll.output) {
			return event
		}
	}
	return EmptyEvent
}
