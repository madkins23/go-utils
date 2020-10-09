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
