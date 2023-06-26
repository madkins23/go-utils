package app

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleSignals(t *testing.T) {
	var userSignals = []os.Signal{syscall.SIGUSR1, syscall.SIGUSR2}
	testSignalHandler(t, func(handler SignalHandler) {
		HandleSignals(handler, userSignals...)
	}, userSignals...)
}

func TestHandleTerminalSignals(t *testing.T) {
	testSignalHandler(t, HandleTerminalSignals, terminalSignals...)
}

func testSignalHandler(
	t *testing.T, testFunc func(SignalHandler), signals ...os.Signal) {
	//
	proc, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	require.NotNil(t, proc)
	for _, sig := range signals {
		ready := make(chan bool)
		require.NotNil(t, ready)
		var signalReceived os.Signal
		testFunc(func(sig os.Signal) {
			signalReceived = sig
			ready <- true
		})
		time.Sleep(5 * time.Millisecond)
		require.NoError(t, proc.Signal(sig))
		_ = <-ready // wait for signal handler to fire
		assert.Equal(t, sig, signalReceived)
	}
}
