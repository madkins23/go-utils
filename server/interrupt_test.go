package server

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterrupt(t *testing.T) {
	ctxt, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	require.NotNil(t, ctxt)
	require.NotNil(t, stop)

	var err error
	var interrupted bool
	go func() {
		time.Sleep(10 * time.Millisecond)
		err = Interrupt()
		interrupted = true
	}()

	select {
	case <-ctxt.Done():
		stop()
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "")
	}

	assert.NoError(t, err)
	assert.True(t, interrupted)
}
