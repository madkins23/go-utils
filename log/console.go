package log

import (
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
