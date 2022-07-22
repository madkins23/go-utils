package error

// ErrorWithDetailMap interface provides a way to determine if an error has string details.
type ErrorWithDetailMap interface {
	error
	DetailStringMap() map[string]string
}

// errWithDetailMap combines an error with an extra map of strings representing error detail.
type errWithDetailMap struct {
	msg    string
	detail map[string]string
}

// NewErrorWithStringMap constructs an errWithStringMap object.
func NewErrorWithStringMap(msg string, detail map[string]string) ErrorWithDetailMap {
	return &errWithDetailMap{
		msg:    msg,
		detail: detail,
	}
}

// NewErrorWithStringArrayDummy provides an empty error object for use with errors.As()
func NewErrorWithStringMapDummy() ErrorWithDetailMap {
	return NewErrorWithStringMap("", nil)
}

// Error() returns the error message without any of the error detail strings.
func (e *errWithDetailMap) Error() string {
	return e.msg
}

// DetailStringMap returns the map of strings containing error detail.
func (e *errWithDetailMap) DetailStringMap() map[string]string {
	return e.detail
}
