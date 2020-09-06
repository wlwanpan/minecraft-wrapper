package main

import (
	"errors"
	"regexp"
)

const (
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelERROR = "ERROR"
)

const (
	logSubmatchCount int = 4
)

var (
	logRegex = regexp.MustCompile(`(\[[0-9:]*\]) \[([A-z #0-9]*)\/([A-z #]*)\]: (.*)`)

	ErrMatchingLog = errors.New("err matching log line")
)

type LogLine struct {
	timestamp  string
	threadName string
	level      string
	output     string
}

func strToLogLine(line string) (*LogLine, error) {
	matches := logRegex.FindAllStringSubmatch(line, logSubmatchCount)
	if len(matches) < 1 {
		return nil, ErrMatchingLog
	}
	subgroups := matches[0]
	return &LogLine{
		timestamp:  subgroups[1],
		threadName: subgroups[2],
		level:      subgroups[3],
		output:     subgroups[4],
	}, nil
}
