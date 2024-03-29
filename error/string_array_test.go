package error

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const stringArrayMsg = "string array test message"

var stringArrayEmpty []string = nil
var stringArrayDetails = []string{
	"alpha",
	"bravo",
	"charlie",
}

func TestErrorWithStringArrayDetails(t *testing.T) {
	err := NewErrorWithStringArray(stringArrayMsg, stringArrayDetails)
	require.Error(t, err)
	assert.Equal(t, stringArrayMsg, err.Error())
	errDetails, ok := err.(WithDetailArray)
	require.True(t, ok)
	assert.Equal(t, stringArrayDetails, errDetails.DetailStringArray())
}

func TestErrorWithNullStringArrayDetails(t *testing.T) {
	err := NewErrorWithStringArray(stringArrayMsg, stringArrayEmpty)
	require.Error(t, err)
	assert.Equal(t, stringArrayMsg, err.Error())
	errDetails, ok := err.(WithDetailArray)
	require.True(t, ok)
	assert.Equal(t, stringArrayEmpty, errDetails.DetailStringArray())

}

func TestErrorAsStringArrayDetails(t *testing.T) {
	err1 := NewErrorWithStringArray(stringArrayMsg, stringArrayDetails)
	require.Error(t, err1)
	assert.Equal(t, stringArrayMsg, err1.Error())
	err2 := fmt.Errorf("wrapped: %w", err1)
	require.Error(t, err2)
	dummy := NewErrorWithStringArrayDummy()
	require.Error(t, dummy)
	require.True(t, errors.As(err2, &dummy))
	assert.IsType(t, NewErrorWithStringArrayDummy(), dummy)
	assert.Equal(t, stringArrayMsg, dummy.Error())
	errDetails, ok := dummy.(WithDetailArray)
	require.True(t, ok)
	assert.Equal(t, stringArrayDetails, errDetails.DetailStringArray())
}

func Example_stringArray() {
	details := make([]string, 3)
	details[0] = "alpha"
	details[1] = "bravo"
	details[2] = "charlie"
	err := NewErrorWithStringArray("message", details)
	wrapped := fmt.Errorf("Error:   %w", err)
	dummy := NewErrorWithStringArrayDummy()
	if errors.As(wrapped, &dummy) {
		details := ""
		if withDetails, ok := err.(WithDetailArray); ok {
			for _, det := range withDetails.DetailStringArray() {
				if len(details) > 0 {
					details += ", "
				}
				details += det
			}
		}
		fmt.Printf("%s\nDetails: %s\n", wrapped, details)
	}

	// Output:
	// Error:   message
	// Details: alpha, bravo, charlie
}
