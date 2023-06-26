package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminator_Shutdown(t *testing.T) {
	term := NewTerminator()
	alpha := &Sub1{}
	assert.False(t, alpha.done)
	term.Add(alpha)
	var bravo Sub2 = "Help!"
	term.Add(bravo)
	charlie := &Sub1{}
	assert.False(t, charlie.done)
	term.Add(charlie)
	var delta Sub2 = "Aaaah!"
	term.Add(delta)
	err := term.Shutdown()
	assert.True(t, alpha.done)
	assert.True(t, charlie.done)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Help!\nAaaah!")
}

type Sub1 struct {
	done bool
}

func (s1 *Sub1) Shutdown() error {
	s1.done = true
	return nil
}

type Sub2 string

func (s2 Sub2) Shutdown() error {
	return errors.New(string(s2))
}
