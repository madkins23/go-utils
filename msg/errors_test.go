package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const name = "name"

func TestErrBlocked(t *testing.T) {
	assert.Equal(t, strBlocked, (&ErrBlocked{}).Error())
	assert.Equal(t, name+strBlockedNamed, (&ErrBlocked{name}).Error())
}

func TestErrDeprecated(t *testing.T) {
	assert.Equal(t, strDeprecated, (&ErrDeprecated{}).Error())
	assert.Equal(t, name+strDeprecatedNamed, (&ErrDeprecated{name}).Error())
}

func TestErrNotImplemented(t *testing.T) {
	assert.Equal(t, strNotImplementedYet, (&ErrNotImplemented{}).Error())
	assert.Equal(t, name+strNotImplementedNamed, (&ErrNotImplemented{name}).Error())
}
