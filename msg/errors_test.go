package msg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	msg  = "eMsg"
	name = "name"
)

func TestConstError_Error(t *testing.T) {
	const e = ConstError(msg)
	assert.Equal(t, msg, e.Error())
	assert.ErrorIs(t, ConstError(msg), e)
}

func TestErrBlocked(t *testing.T) {
	assert.Equal(t, strBlocked, (&ErrBlocked{}).Error())
	assert.Equal(t, name+strBlockedNamed, (&ErrBlocked{name}).Error())
	assert.ErrorIs(t, &ErrBlocked{}, &ErrBlocked{})
	assert.ErrorIs(t, &ErrBlocked{Name: "test"}, &ErrBlocked{Name: "test"}, "exact match")
	assert.ErrorIs(t, &ErrBlocked{Name: "test"}, &ErrBlocked{}, "generic match")
	assert.False(t, errors.Is(&ErrBlocked{}, &ErrBlocked{Name: "test"}), "specific fail")
	assert.False(t, errors.Is(&ErrBlocked{Name: "fail"}, &ErrBlocked{Name: "test"}), "mismatch fail")
}

func TestErrDeprecated(t *testing.T) {
	assert.Equal(t, strDeprecated, (&ErrDeprecated{}).Error())
	assert.Equal(t, name+strDeprecatedNamed, (&ErrDeprecated{name}).Error())
	assert.ErrorIs(t, &ErrDeprecated{}, &ErrDeprecated{})
	assert.ErrorIs(t, &ErrDeprecated{Name: "test"}, &ErrDeprecated{Name: "test"}, "exact match")
	assert.ErrorIs(t, &ErrDeprecated{Name: "test"}, &ErrDeprecated{}, "generic match")
	assert.False(t, errors.Is(&ErrDeprecated{}, &ErrDeprecated{Name: "test"}), "specific fail")
	assert.False(t, errors.Is(&ErrDeprecated{Name: "fail"}, &ErrDeprecated{Name: "test"}), "mismatch fail")
}

func TestErrNotImplemented(t *testing.T) {
	assert.Equal(t, strNotImplementedYet, (&ErrNotImplemented{}).Error())
	assert.Equal(t, name+strNotImplementedNamed, (&ErrNotImplemented{name}).Error())
	assert.ErrorIs(t, &ErrNotImplemented{}, &ErrNotImplemented{})
	assert.ErrorIs(t, &ErrNotImplemented{Name: "test"}, &ErrNotImplemented{Name: "test"}, "exact match")
	assert.ErrorIs(t, &ErrNotImplemented{Name: "test"}, &ErrNotImplemented{}, "generic match")
	assert.False(t, errors.Is(&ErrNotImplemented{}, &ErrNotImplemented{Name: "test"}), "specific fail")
	assert.False(t, errors.Is(&ErrNotImplemented{Name: "fail"}, &ErrNotImplemented{Name: "test"}), "mismatch fail")
}
