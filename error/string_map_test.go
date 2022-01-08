package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const stringMapMsg = "string map test message"

var stringMapEmpty map[string]string = nil
var stringMapDetails = map[string]string{
	"alpha":   "one",
	"bravo":   "two",
	"charlie": "three",
}

func TestErrorWithStringMapDetails(t *testing.T) {
	err := NewErrorWithStringMap(stringMapMsg, stringMapDetails)
	require.Error(t, err)
	assert.Equal(t, stringMapMsg, err.Error())
	errDetails, ok := err.(StringMapDetails)
	require.True(t, ok)
	assert.Equal(t, stringMapDetails, errDetails.DetailStringMap())
}

func TestErrorWithNullStringMapDetails(t *testing.T) {
	err := NewErrorWithStringMap(stringMapMsg, stringMapEmpty)
	require.Error(t, err)
	assert.Equal(t, stringMapMsg, err.Error())
	errDetails, ok := err.(StringMapDetails)
	require.True(t, ok)
	assert.Equal(t, stringMapEmpty, errDetails.DetailStringMap())
}

func Example_stringMap() {
	details := make(map[string]string, 3)
	details["one"] = "alpha"
	details["two"] = "bravo"
	details["three"] = "charlie"
	err := NewErrorWithStringMap("message", details)
	fmt.Printf("Error: %s\n", err)
	if withDetails, ok := err.(StringMapDetails); ok {
		for _, det := range withDetails.DetailStringMap() {
			fmt.Printf("       %s\n", det)
		}
	}

	// Output:
	// Error: message
	//        alpha
	//        bravo
	//        charlie
}
