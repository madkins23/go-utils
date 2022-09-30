package log

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Console configures logging to output to console in a more readable format.
// This will be slower logging than just using basic JSON output.
// Output is sent to os.Stderr and timestamps are in local time.
func Console() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
}

// ConsoleOrFile is used to open a console log or a log file based on its fields.
// If a log file is opened it may be specified as JSON objects
// or the more readable console log format.
// The fields should be set based on command line arguments.
type ConsoleOrFile struct {
	Console bool
	LogFile string
	AsJSON  bool
}

// Setup opens the console log or log file as appropriate based on the object's fields.
func (cof *ConsoleOrFile) Setup() error {
	if cof.Console {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	} else if f, err := os.OpenFile(cof.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		return fmt.Errorf("Log file creation: %w", err)
	} else {
		defer func() { _ = f.Close() }()
		if cof.AsJSON {
			log.Logger = log.Output(f)
		} else {
			// Separate blocks of log statements for each run.
			_, _ = fmt.Fprintln(f)
			// Use ConsoleWriter for readable text instead of JSON blocks.
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: f, TimeFormat: "15:04:05", NoColor: true})
		}
	}
	return nil
}
