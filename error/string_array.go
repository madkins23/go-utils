package error

// ErrorWithDetailArray interface provides a way to determine if an error has string details.
type ErrorWithDetailArray interface {
	error
	DetailStringArray() []string
}

// errWithDetailArray combines an error with an extra array of strings representing error detail.
type errWithDetailArray struct {
	msg    string
	detail []string
}

// NewErrorWithStringArray constructs an errWithStringArray object.
func NewErrorWithStringArray(msg string, detail []string) ErrorWithDetailArray {
	return &errWithDetailArray{
		msg:    msg,
		detail: detail,
	}
}

// NewErrorWithStringArrayDummy provides an empty error object for use with errors.As()
func NewErrorWithStringArrayDummy() ErrorWithDetailArray {
	return NewErrorWithStringArray("", nil)
}

// Error() returns the error message without any of the error detail strings.
func (e *errWithDetailArray) Error() string {
	return e.msg
}

// DetailStringArray returns the array of strings containing error detail.
func (e *errWithDetailArray) DetailStringArray() []string {
	return e.detail
}
