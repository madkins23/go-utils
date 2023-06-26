package app

import (
	"errors"
	"sync"
)

// SubSystem to be shut down by Terminator.
type SubSystem interface {
	Shutdown() error
}

// Terminator is used to shut down subsystems gracefully.
type Terminator struct {
	lock       sync.Mutex
	shutDown   bool
	subSystems []SubSystem
}

func NewTerminator() *Terminator {
	return &Terminator{
		subSystems: make([]SubSystem, 0),
	}
}

func (t *Terminator) Add(subSystem SubSystem) {
	t.subSystems = append(t.subSystems, subSystem)
}

func (t *Terminator) Shutdown() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.shutDown {
		return nil
	}

	errs := make([]error, len(t.subSystems))
	for i, subSystem := range t.subSystems {
		errs[i] = subSystem.Shutdown()
	}
	t.shutDown = true
	return errors.Join(errs...)
}
