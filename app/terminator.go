package app

import (
	"errors"
)

// SubSystem to be shut down by Terminator.
type SubSystem interface {
	Shutdown() error
}

// Terminator is used to shut down subsystems gracefully.
type Terminator struct {
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
	errs := make([]error, len(t.subSystems))
	for i, subSystem := range t.subSystems {
		errs[i] = subSystem.Shutdown()
	}
	return errors.Join(errs...)
}
