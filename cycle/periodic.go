package cycle

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// PeriodicFn is a function to be executed periodically.
type PeriodicFn func(uint) error

// FinalFn is a function to be executed after periodic behavior has terminated.
type FinalFn func()

// SignalFn is a function to be executed when a signal is received.
// Calling code may use this for logging or debugging.
type SignalFn func(os.Signal)

// Periodic function execution object.
type Periodic struct {
	done     chan bool
	stop     chan bool
	signals  chan os.Signal
	cycleFn  PeriodicFn
	finalFn  FinalFn
	signalFn SignalFn
	endErr   error
}

var errNoCycleFn = errors.New("no cycle function")

// NewPeriodic returns a new Periodic object or an error.
// The cycleFn argument is required, the other two can be nil.
func NewPeriodic(cycleFn PeriodicFn, finalFn FinalFn, signalFn SignalFn) (*Periodic, error) {
	if cycleFn == nil {
		return nil, errNoCycleFn
	}

	p := &Periodic{
		stop:     make(chan bool, 3),
		done:     make(chan bool, 2),
		signals:  make(chan os.Signal),
		cycleFn:  cycleFn,
		finalFn:  finalFn,
		signalFn: signalFn,
	}

	return p, nil
}

// handleSignals is run as a goroutine to handle termination signals (e.g. <ctrl>-C).
func (p *Periodic) handleSignals() {
	if sig, ok := <-p.signals; ok {
		if p.signalFn != nil {
			p.signalFn(sig)
		}
		p.stop <- true
	}
}

// Execute cycle function on a regular interval.
// Uses a ticker to start the code at known intervals.
// If the code runs longer than a ticker interval some intervals will be skipped.
// The code is run initially, then the ticker is started.
// Run in a goroutine or this method will block until completion.
func (p *Periodic) Ticker(interval time.Duration) {
	if p.finalFn != nil {
		defer p.finalFn()
	}

	cycles := uint(0)
	if p.endErr = p.cycleFn(cycles); p.endErr != nil {
		p.done <- true
		return
	}

	signal.Notify(p.signals, syscall.SIGINT, syscall.SIGTERM)
	go p.handleSignals()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		// Use two select statements to prioritize stop channel over ticker.
		select {
		case <-p.stop:
			close(p.signals)
			p.done <- true
			ticker.Stop()
			return
		default:
		}
		select {
		case <-ticker.C:
			cycles++
		}
		if p.endErr = p.cycleFn(cycles); p.endErr != nil {
			p.stop <- true
			continue
		}
	}
}

// Stop periodic cycling.
func (p *Periodic) Stop() {
	p.stop <- true
}

// Finished waits for periodic cycling to end and returns any final error.
// Does not call Stop().
func (p *Periodic) Finished() error {
	<-p.done
	return p.endErr
}
