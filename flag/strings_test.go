package flag

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringArrayFlagsSingle(t *testing.T) {
	os.Args = []string{"TestStringArrayFlags", "-cmdSingle", "one"}
	var testFlags StringArray
	flag.Var(&testFlags, "cmdSingle", "Some description for this param.")
	flag.Parse()
	require.NotNil(t, testFlags)
	assert.Len(t, testFlags, 1)
	assert.Equal(t, "one", testFlags[0])
}

func TestStringArrayFlagsMultiple(t *testing.T) {
	os.Args = []string{"TestStringArrayFlags", "-cmdMultiple", "one", "-cmdMultiple", "two"}
	var testFlags StringArray
	flag.Var(&testFlags, "cmdMultiple", "Some description for this param.")
	flag.Parse()
	require.NotNil(t, testFlags)
	assert.Len(t, testFlags, 2)
	assert.Equal(t, "one", testFlags[0])
	assert.Equal(t, "two", testFlags[1])
}

func TestStringArrayFlagsCommas(t *testing.T) {
	os.Args = []string{"TestStringArrayFlags", "-cmdCommas", "one,two,three"}
	var testFlags StringArray
	flag.Var(&testFlags, "cmdCommas", "Some description for this param.")
	flag.Parse()
	require.NotNil(t, testFlags)
	assert.Len(t, testFlags, 3)
	assert.Equal(t, "one", testFlags[0])
	assert.Equal(t, "two", testFlags[1])
	assert.Equal(t, "three", testFlags[2])
}

func Example_stringArrayFlags() {
	os.Args = []string{"TestStringArrayFlags", "-cmd", "one", "-cmd", "two,three", "-cmd", "four"}
	var testFlags StringArray
	flag.Var(&testFlags, "cmd", "Some description for this param.")
	flag.Parse()
	fmt.Println(testFlags)
	// Output: [one two three four]
}
