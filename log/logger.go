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

// Debug returns a debug event from the default zerolog logger.
func Debug() *zerolog.Event {
	return log.Logger.Debug()
}

// Error returns a error event from the default zerolog logger.
func Error() *zerolog.Event {
	return log.Logger.Error()
}

// Fatal returns a fatal event from the default zerolog logger.
func Fatal() *zerolog.Event {
	return log.Logger.Fatal()
}

// Info returns a info event from the default zerolog logger.
func Info() *zerolog.Event {
	return log.Logger.Info()
}

// Panic returns a panic event from the default zerolog logger.
func Panic() *zerolog.Event {
	return log.Logger.Info()
}

// Trace returns a trace event from the default zerolog logger.
func Trace() *zerolog.Event {
	return log.Logger.Trace()
}

// Warn returns a warn event from the default zerolog logger.
func Warn() *zerolog.Event {
	return log.Logger.Warn()
}
