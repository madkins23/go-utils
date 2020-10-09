package test

import (
	"bytes"
	"io"
	"os"
)

// CaptureStderr executes the specified function,
// capturing all stderr output and returning it along with any error result.
func CaptureStderr(test func()) (string, error) {
	return capture(&os.Stderr, test)
}

// CaptureStdout executes the specified function,
// capturing all stdout output and returning it along with any error result.
func CaptureStdout(test func()) (string, error) {
	return capture(&os.Stdout, test)
}

func capture(orig **os.File, test func()) (string, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = writer.Close()
	}()

	original := *orig
	*orig = writer
	defer func() {
		*orig = original
	}()

	test()
	err = writer.Close()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	return buf.String(), err
}
