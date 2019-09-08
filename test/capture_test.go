package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testErrString = "Standard Error"
	testOutString = "Standard Output"
)

func TestCaptureStderr(t *testing.T) {
	text, err := CaptureStderr(func() {
		_, _ = fmt.Fprint(os.Stderr, testErrString)
	})
	assert.NoError(t, err)
	assert.Equal(t, testErrString, text)
}

func TestCaptureStdout(t *testing.T) {
	text, err := CaptureStdout(func() {
		_, _ = fmt.Fprint(os.Stdout, testOutString)
	})
	assert.NoError(t, err)
	assert.Equal(t, testOutString, text)
}

func TestCaptureBoth(t *testing.T) {
	var textErr string
	textOut, err := CaptureStdout(func() {
		var err error
		textErr, err = CaptureStderr(func() {
			_, _ = fmt.Fprint(os.Stderr, testErrString)
			_, _ = fmt.Fprint(os.Stdout, testOutString)
		})
		assert.NoError(t, err)
	})
	assert.NoError(t, err)
	assert.Equal(t, testErrString, textErr)
	assert.Equal(t, testOutString, textOut)
}
