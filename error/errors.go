package error

import "errors"

const (
	deprecated      = "deprecated: "
	toBeImplemented = "to be implemented: "
)

var notYetImplemented = errors.New("not yet implemented")

func NotYetImplemented() error {
	return notYetImplemented
}

func ToBeImplemented(name string) error {
	return errors.New(toBeImplemented + name)
}

func Deprecated(name string) error {
	return errors.New(deprecated + name)
}
