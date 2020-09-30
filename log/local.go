package log

import (
	"github.com/rs/zerolog"
)

// LocalLogger provides an object-specific logger and associated functionality.
type LocalLogger struct {
	logger *zerolog.Logger
}

// Logger returns the local logger.
// If the local logger is not yet set it will be configured to the default logger.
// Override this method to define another logger.
func (ll *LocalLogger) Logger() *zerolog.Logger {
	return ll.LoggerWithFn(Logger)
}

// LocalLogger returns the local logger.
// If the local logger is not yet set it will be configured to the logger returned
// by the specified default function.
func (ll *LocalLogger) LoggerWithFn(defaultFn func() *zerolog.Logger) *zerolog.Logger {
	if ll.logger == nil {
		ll.logger = defaultFn()
	}

	return ll.logger
}

// SetLogger sets the local logger.
func (ll *LocalLogger) SetLogger(logger *zerolog.Logger) {
	ll.logger = logger
}
