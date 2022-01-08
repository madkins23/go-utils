package error

import (
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
	errDetails, ok := err.(StringArrayDetails)
	require.True(t, ok)
	assert.Equal(t, stringArrayDetails, errDetails.DetailStringArray())
}

func TestErrorWithNullStringArrayDetails(t *testing.T) {
	err := NewErrorWithStringArray(stringArrayMsg, stringArrayEmpty)
	require.Error(t, err)
	assert.Equal(t, stringArrayMsg, err.Error())
	errDetails, ok := err.(StringArrayDetails)
	require.True(t, ok)
	assert.Equal(t, stringArrayEmpty, errDetails.DetailStringArray())
}

func Example_stringArray() {
	details := make([]string, 3)
	details[0] = "alpha"
	details[1] = "bravo"
	details[2] = "charlie"
	err := NewErrorWithStringArray("message", details)
	fmt.Printf("Error: %s\n", err)
	if withDetails, ok := err.(StringArrayDetails); ok {
		for _, det := range withDetails.DetailStringArray() {
			fmt.Printf("       %s\n", det)
		}
	}

	// Output:
	// Error: message
	//        alpha
	//        bravo
	//        charlie
}
