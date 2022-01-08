package error

// StringArrayDetails interface provides a way to determine if an error has string details.
type StringArrayDetails interface {
	DetailStringArray() []string
}

// errWithStringArray combines an error with an extra array of strings representing error detail.
type errWithStringArray struct {
	msg    string
	detail []string
}

// NewErrorWithStringArray constructs an errWithStringArray object.
func NewErrorWithStringArray(msg string, detail []string) error {
	return &errWithStringArray{
		msg:    msg,
		detail: detail,
	}
}

// Error() returns the error message without any of the error detail strings.
func (e *errWithStringArray) Error() string {
	return e.msg
}

// DetailStringArray returns the array of strings containing error detail.
func (e *errWithStringArray) DetailStringArray() []string {
	return e.detail
}
