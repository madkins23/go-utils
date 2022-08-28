package error

import "errors"

var notYetImplemented = errors.New("not yet implemented")

// NotYetImplemented returns a generic TBD error.
func NotYetImplemented() error {
	return notYetImplemented
}

const (
	deprecated      = "deprecated: "
	toBeImplemented = "to be implemented: "
)

// ToBeImplemented returns a TBD error specifying the name of the unimplemented function or method.
func ToBeImplemented(name string) error {
	return errors.New(toBeImplemented + name)
}

// Deprecated returns a deprecation error specifying the name of the deprecated function or method.
func Deprecated(name string) error {
	return errors.New(deprecated + name)
}
