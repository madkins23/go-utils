package cycle

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/madkins23/go-utils/log"
	"github.com/rs/zerolog"
)

// CycleFn is a function to be executed cyclically.
type CycleFn func(uint) error

// FinalFn is a function to be executed after cyclic behavior has terminated.
type FinalFn func()

// Periodic function execution.
type Periodic struct {
	log.LocalLogger
	done    chan bool
	stop    chan bool
	signals chan os.Signal
	cycleFn CycleFn
	finalFn FinalFn
	endErr  error
}

var errNoCycleFn = errors.New("no cycle function")

func NewPeriodic(logger *zerolog.Logger, cycleFn CycleFn, finalFn FinalFn) (*Periodic, error) {
	if cycleFn == nil {
		return nil, errNoCycleFn
	}

	p := &Periodic{
		stop:    make(chan bool, 3),
		done:    make(chan bool, 2),
		signals: make(chan os.Signal),
		cycleFn: cycleFn,
		finalFn: finalFn,
	}
	p.SetLogger(logger)

	return p, nil
}

func (p *Periodic) handleSignals() {
	if sig, ok := <-p.signals; ok {
		p.Logger().Info().Str("signal", sig.String()).Msg("Received Signal")
		p.stop <- true
	}
}

// Execute cycle function on a regular interval.
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

// Wait for application to end and return final error.
// Does not call Stop().
func (p *Periodic) Finished() error {
	<-p.done
	return p.endErr
}
