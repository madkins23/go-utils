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
