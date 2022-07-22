package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsZero(t *testing.T) {
	assert.True(t, IsZero(0))
	assert.False(t, IsZero(1))
	assert.True(t, IsZero(false))
	assert.False(t, IsZero(true))
}

func TestErrorIfZero(t *testing.T) {
	assert.EqualError(t, ErrorIfZero(0), ErrIsZero.Error())
	assert.NoError(t, ErrorIfZero(1))
	assert.EqualError(t, ErrorIfZero(false), ErrIsZero.Error())
	assert.NoError(t, ErrorIfZero(true))
}

func TestIsZero_Pointer(t *testing.T) {
	var x *int
	assert.True(t, x == nil)
	assert.True(t, IsZero(x))
	assert.EqualError(t, ErrorIfZero(x), ErrIsZero.Error())
	x = new(int)
	assert.False(t, IsZero(x))
	assert.NoError(t, ErrorIfZero(x))
}

// This is the case that resulted in check.IsZero():

type something interface {
}

type container[T something] struct {
	s T
}

func (c *container[T]) do() error {
	return ErrorIfZero(c.s)
}

func TestIsZero_Generic(t *testing.T) {
	c := new(container[something])
	assert.NotNil(t, c)
	assert.Nil(t, c.s)
	assert.EqualError(t, c.do(), ErrIsZero.Error())
}
