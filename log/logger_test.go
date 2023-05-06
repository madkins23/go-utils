package log

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleLogger() {
	Logger().Info().Msg("Message")
	// Output:
}

func TestLogger(t *testing.T) {
	logger := Logger()
	require.NotNil(t, logger)
	assert.IsType(t, &zerolog.Logger{}, logger)
}

func TestSetLogger(t *testing.T) {
	oldLogger := Logger()
	require.NotNil(t, oldLogger)
	defer SetLogger(*oldLogger)
	newLogger := Logger().With().Str("key", "value").Logger()
	require.NotEqual(t, Logger(), &newLogger)
	SetLogger(newLogger)
	require.NotNil(t, Logger())
	require.Equal(t, Logger(), &newLogger)
	assert.IsType(t, &zerolog.Logger{}, Logger())
}
