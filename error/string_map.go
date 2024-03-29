package error

// WithDetailMap interface provides a way to determine if an error has string details.
//
// Deprecated: No obvious replacement.
type WithDetailMap interface {
	error
	DetailStringMap() map[string]string
}

// ErrorWithDetailMap interface provides a way to determine if an error has string details.
//
// Deprecated: No obvious replacement.
type ErrorWithDetailMap interface {
	error
	DetailStringMap() map[string]string
}

// withDetailMap combines an error with an extra map of strings representing error detail.
type withDetailMap struct {
	msg    string
	detail map[string]string
}

// NewErrorWithStringMap constructs an error with a map of strings representing error details.
//
// Deprecated: No obvious replacement.
func NewErrorWithStringMap(msg string, detail map[string]string) WithDetailMap {
	return &withDetailMap{
		msg:    msg,
		detail: detail,
	}
}

// NewErrorWithStringMapDummy provides an empty error object for use with errors.As()
//
// Deprecated: No obvious replacement.
func NewErrorWithStringMapDummy() WithDetailMap {
	return NewErrorWithStringMap("", nil)
}

// Error() returns the error message without any of the error detail strings.
func (e *withDetailMap) Error() string {
	return e.msg
}

// DetailStringMap returns the map of strings containing error detail.
func (e *withDetailMap) DetailStringMap() map[string]string {
	return e.detail
}
