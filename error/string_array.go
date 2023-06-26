package error

// WithDetailArray interface provides a way to determine if an error has string details.
//
// Deprecated: Use Go package errors.Join() instead.
type WithDetailArray interface {
	error
	DetailStringArray() []string
}

// ErrorWithDetailArray interface provides a way to determine if an error has string details.
//
// Deprecated: Name begins with package name, use error.WithDetailArray instead.
type ErrorWithDetailArray interface {
	error
	DetailStringArray() []string
}

// withDetailArray combines an error with an array of strings representing error details.
type withDetailArray struct {
	msg    string
	detail []string
}

// NewErrorWithStringArray constructs an error with an array of string details.
//
// Deprecated: Use Go package errors.Join() instead.
func NewErrorWithStringArray(msg string, detail []string) WithDetailArray {
	return &withDetailArray{
		msg:    msg,
		detail: detail,
	}
}

// NewErrorWithStringArrayDummy provides an empty error object for use with errors.As()
//
// Deprecated: Use Go package errors.Join() instead.
func NewErrorWithStringArrayDummy() WithDetailArray {
	return NewErrorWithStringArray("", nil)
}

// Error() returns the error message without any of the error detail strings.
func (e *withDetailArray) Error() string {
	return e.msg
}

// DetailStringArray returns the array of strings containing error detail.
func (e *withDetailArray) DetailStringArray() []string {
	return e.detail
}
