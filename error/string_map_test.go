package error

import (
	"errors"
	"fmt"
	"sort"
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

func TestErrorAsStringMapDetails(t *testing.T) {
	err1 := NewErrorWithStringMap(stringMapMsg, stringMapDetails)
	require.Error(t, err1)
	assert.Equal(t, stringMapMsg, err1.Error())
	err2 := fmt.Errorf("wrapped: %w", err1)
	require.Error(t, err2)
	dummy := NewErrorWithStringMap("", nil)
	require.Error(t, dummy)
	require.True(t, errors.As(err2, &dummy))
	assert.IsType(t, NewErrorWithStringMapDummy(), dummy)
	assert.Equal(t, stringMapMsg, dummy.Error())
	errDetails, ok := dummy.(StringMapDetails)
	require.True(t, ok)
	assert.Equal(t, stringMapDetails, errDetails.DetailStringMap())
}

func Example_stringMap() {
	details := make(map[string]string, 3)
	details["1"] = "alpha"
	details["2"] = "bravo"
	details["3"] = "charlie"
	err := NewErrorWithStringMap("message", details)
	wrapped := fmt.Errorf("Wrapped: %w", err)
	fmt.Printf("Error: %s\n", wrapped)
	dummy := NewErrorWithStringMapDummy()
	if errors.As(wrapped, &dummy) {
		if withDetails, ok := err.(StringMapDetails); ok {
			detailed := withDetails.DetailStringMap()
			// Maps return keys/values in a deliberately random order,
			// so we must get the keys, sort them, and then use them.
			keys := make([]string, 0, len(detailed))
			for k, _ := range detailed {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, key := range keys {
				fmt.Printf("       %s\n", detailed[key])
			}
		}
	}

	// Output:
	// Error: Wrapped: message
	//        alpha
	//        bravo
	//        charlie
}
