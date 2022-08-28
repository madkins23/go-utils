package error

import "errors"

var notYetImplemented = errors.New("not yet implemented")

// NotYetImplemented returns a generic TBD error.
//
// Deprecated: use struct msg.ErrNotImplemented.
func NotYetImplemented() error {
	return notYetImplemented
}

const (
	deprecated      = "deprecated: "
	toBeImplemented = "to be implemented: "
)

// ToBeImplemented returns a TBD error specifying the name of the unimplemented function or method.
//
// Deprecated: use struct msg.ErrNotImplemented.
func ToBeImplemented(name string) error {
	return errors.New(toBeImplemented + name)
}

// Deprecated returns a deprecation error specifying the name of the deprecated function or method.
//
// Deprecated: use struct msg.ErrDeprecated.
func Deprecated(name string) error {
	return errors.New(deprecated + name)
}
