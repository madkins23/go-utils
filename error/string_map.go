package error

// StringMapDetails interface provides a way to determine if an error has string details.
type StringMapDetails interface {
	error
	DetailStringMap() map[string]string
}

// errWithStringMap combines an error with an extra map of strings representing error detail.
type errWithStringMap struct {
	msg    string
	detail map[string]string
}

// NewErrorWithStringMap constructs an errWithStringMap object.
func NewErrorWithStringMap(msg string, detail map[string]string) StringMapDetails {
	return &errWithStringMap{
		msg:    msg,
		detail: detail,
	}
}

// NewErrorWithStringArrayDummy provides an empty error object for use with errors.As()
func NewErrorWithStringMapDummy() StringMapDetails {
	return NewErrorWithStringMap("", nil)
}

// Error() returns the error message without any of the error detail strings.
func (e *errWithStringMap) Error() string {
	return e.msg
}

// DetailStringMap returns the map of strings containing error detail.
func (e *errWithStringMap) DetailStringMap() map[string]string {
	return e.detail
}
