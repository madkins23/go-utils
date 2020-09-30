package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger returns a default logger.
func Logger() *zerolog.Logger {
	return &log.Logger
}
