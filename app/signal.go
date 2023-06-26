package app

import (
	"os"
	"os/signal"
	"syscall"
)

type SignalHandler func(sig os.Signal)

func HandleSignals(handler SignalHandler, signals ...os.Signal) {
	channel := make(chan os.Signal)
	go func() {
		handler(<-channel)
	}()
	signal.Notify(channel, signals...)
}

// terminalSignals is the list of signals that might normally terminate a program.
// Go can't handle SIGKILL or SIGSTOP:
// https://pkg.go.dev/os/signal#hdr-Types_of_signals
var terminalSignals = []os.Signal{
	syscall.SIGABRT,
	syscall.SIGBUS,
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTERM,
}

func HandleTerminalSignals(handler SignalHandler) {
	HandleSignals(handler, terminalSignals...)
}
