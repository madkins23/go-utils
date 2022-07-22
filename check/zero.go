package check

import (
	"errors"
	"reflect"
)

var ErrIsZero = errors.New("item is zero")

// IsZero checks to see if the specified entity is its zero value.
// This is particularly problematic in interface and/or generic contexts.
// May be useful if you are getting this compile error:
//  invalid operation: c.s == nil (mismatched types T and untyped nil)
// Too bad this requires reflection.  :-(
func IsZero[T any](x T) bool {
	return reflect.ValueOf(&x).Elem().IsZero()
}

// ErrorIfZero returns the error ErrIsZero if the specified entity is its zero value.
func ErrorIfZero[T any](x T) error {
	if IsZero(x) {
		return ErrIsZero
	}
	return nil
}
