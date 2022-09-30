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
	logFile *os.File
}

// Setup opens the console log or log file as appropriate based on the object's fields.
func (cof *ConsoleOrFile) Setup() error {
	var err error
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	if cof.Console {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	} else if cof.logFile, err = os.OpenFile(cof.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		return fmt.Errorf("Log file creation: %w", err)
	} else {
		if cof.AsJSON {
			log.Logger = log.Output(cof.logFile)
		} else {
			// Separate blocks of log statements for each run.
			_, _ = fmt.Fprintln(cof.logFile)
			// Use ConsoleWriter for readable text instead of JSON blocks.
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: cof.logFile, TimeFormat: "15:04:05", NoColor: true})
		}
	}
	return nil
}

// Close any log file that may have been opened.
func (cof *ConsoleOrFile) Close() error {
	if cof.logFile != nil {
		return cof.logFile.Close()
	}
	return nil
}

// CloseForDefer closes any log file that may have been opened.
// Returns no error so this call is simpler to use within a defer statement.
func (cof *ConsoleOrFile) CloseForDefer() {
	_ = cof.Close()
}
