package log

import (
	"flag"
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
// The log to be used is the default zerolog object.
// If a log file is opened it may be specified as JSON objects
// or the more readable console log format.
//
// This struct has been configured with JSON and YAML struct tags.
// This supports reading the contents from a JSON or YAML configuration file.
// The struct can also be embedded within a JSON or YAML configuration file.
type ConsoleOrFile struct {
	// True to send output to console (stderr by default) instead of file.
	Console bool `json:"console" yaml:"console"`

	// True to send console output to stdout instead of stderr
	UseStdout bool `json:"useStdout" yaml:"useStdout"`

	// Path to logfile to be used if not console mode.
	LogFile string `json:"logFile" yaml:"logFile"`

	// True to dump data as JSON to logfile,
	// otherwise console mode output with blank lines when application restarted.
	AsJSON bool `json:"asJson" yaml:"asJson"`

	// ConsoleWriter object available to change settings or nil if sending JSON to logfile.
	writer *zerolog.ConsoleWriter

	// Log file object if used or nil if Console is true.
	logFile *os.File
}

// AddFlagsToSet adds flags to the specified flag.FlagSet.
// Each flag is connected to a field in the configuration object.
func (cof *ConsoleOrFile) AddFlagsToSet(flags *flag.FlagSet, logFile string) {
	flags.BoolVar(&cof.Console, "console", false, "Log to the console instead of the specified log file")
	flags.StringVar(&cof.LogFile, "logFile", logFile, "Log file path")
	flags.BoolVar(&cof.AsJSON, "logJSON", false, "Log output to file as JSON objects")
}

// Setup opens the console log or log file as appropriate based on the object's fields.
func (cof *ConsoleOrFile) Setup() error {
	var err error
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	if cof.Console {
		out := os.Stderr
		if cof.UseStdout {
			out = os.Stdout
		}
		cof.writer = &zerolog.ConsoleWriter{Out: out, TimeFormat: "15:04:05"}
		log.Logger = log.Output(cof.writer)
	} else if cof.logFile, err = os.OpenFile(cof.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		return fmt.Errorf("log file creation: %w", err)
	} else if cof.AsJSON {
		log.Logger = log.Output(cof.logFile)
	} else {
		// Separate blocks of log statements for each run.
		_, _ = fmt.Fprintln(cof.logFile)
		// Use ConsoleWriter for readable text instead of JSON blocks.
		cof.writer = &zerolog.ConsoleWriter{Out: cof.logFile, TimeFormat: "15:04:05", NoColor: true}
		log.Logger = log.Output(cof.writer)
	}
	return nil
}

// Writer returns a pointer to the zerolog.ConsoleWriter created in Setup
// or nil if a log file was configured.
// This object can be used to derive console loggers with additional
// configuration later in their usage.
func (cof *ConsoleOrFile) Writer() *zerolog.ConsoleWriter {
	return cof.writer
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
