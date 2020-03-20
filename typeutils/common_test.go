package typeutils

import "fmt"

type Actor interface {
	declaim() string
}

type alpha struct {
	what string
}

type bravo struct {
	ok bool
}

func (a *alpha) declaim() string {
	return a.what
}

func (b *bravo) declaim() string {
	return fmt.Sprintf("OK:  %t", b.ok)
}
