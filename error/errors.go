package error

import "errors"

var notYetImplemented = errors.New("not yet implemented")

func NotYetImplemented() error {
	return notYetImplemented
}

func ToBeImplemented(name string) error {
	return errors.New("to be implemented: " + name)
}
