package log

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type withLocalLogger struct {
	LocalLogger
}

func ExampleLocalLogger_Logger() {
	type useDefaultLogger struct {
		LocalLogger
	}
	var object useDefaultLogger
	object.Logger().Info().Msg("Use default logger")
	// Output:
}

func ExampleLocalLogger_LoggerWithFn() {
	type useSpecialLogger struct {
		LocalLogger
	}
	var object useSpecialLogger
	// Override LocalLogger.Logger() by defining userSpecialLogger.Logger()
	// and have it call LocalLogger.LoggerWithFn() sort of as shown below
	// (can't define methods within example function):
	object.LoggerWithFn(func() *zerolog.Logger {
		defaultLogger := Logger().With().Str("fn", "special").Logger()
		return &defaultLogger
	}).Info().Msg("Use special logger")
	// Output:
}

func ExampleLocalLogger_SetLogger() {
	type useSetLogger struct {
		LocalLogger
	}
	var object useSetLogger
	setLogger := Logger().With().Str("set", "logger").Logger()
	object.SetLogger(&setLogger)
	object.Logger().Info().Msg("Use set logger")
	// Output:
}

func TestLocalLogger(t *testing.T) {
	w := &withLocalLogger{}
	require.NotNil(t, w)
	assert.Nil(t, w.logger)
	l := w.Logger()
	require.NotNil(t, l)
	assert.IsType(t, &zerolog.Logger{}, l)
	assert.NotNil(t, w.logger)
}

func TestLocalLoggerDefaultFn(t *testing.T) {
	w := &withLocalLogger{}
	require.NotNil(t, w)
	assert.Nil(t, w.logger)
	logger := Logger().With().Str("test", "defaultFn").Logger()
	l := w.LoggerWithFn(func() *zerolog.Logger {
		return &logger
	})
	require.NotNil(t, l)
	assert.IsType(t, &zerolog.Logger{}, l)
	assert.Equal(t, &logger, l)
	assert.NotNil(t, w.logger)
	assert.Equal(t, &logger, w.logger)
}

func TestLocalSetLogger(t *testing.T) {
	w := &withLocalLogger{}
	require.NotNil(t, w)
	assert.Nil(t, w.logger)
	logger := Logger().With().Str("test", "set").Logger()
	w.SetLogger(&logger)
	l := w.Logger()
	require.NotNil(t, l)
	assert.IsType(t, &zerolog.Logger{}, l)
	assert.Equal(t, &logger, l)
	assert.NotNil(t, w.logger)
	assert.Equal(t, &logger, w.logger)
}
