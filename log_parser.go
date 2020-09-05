package main

import "regexp"

const (
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelERROR = "ERROR"
)

var (
	logRegex = regexp.MustCompile(`(\[[0-9:]*\]) \[([A-z #0-9]*)\/([A-z #]*)\]: (.*)`)
)

type Update struct {
	timestamp string // Parse to time.Time
	logLevel  string
}

func logLineToUpdate(line string) *Update {
	matches := logRegex.FindAllStringSubmatch(line, 5)

	// log_time = r.group(1)
	// server_thread = r.group(2)
	// log_level = r.group(3)
	// output = r.group(4)
	// for _, i := range logData {
	// 	log.Println(i)
	// }
	return nil
}
