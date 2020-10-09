package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger returns the default zerolog logger.
// It's kind of redundant but used by other functionality in this package.
func Logger() *zerolog.Logger {
	return &log.Logger
}
