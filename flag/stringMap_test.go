package flag

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapStringFlagsSingle(t *testing.T) {
	os.Args = []string{"TestStringArrayFlags", "-mapSingle", "one:alpha"}
	var testFlags StringMap
	flag.Var(&testFlags, "mapSingle", "Some description for this param.")
	flag.Parse()
	require.NotNil(t, testFlags)
	assert.Len(t, testFlags, 1)
	assert.Equal(t, "alpha", testFlags["one"])
}

func TestMapStringFlagsMultiple(t *testing.T) {
	os.Args = []string{"TestStringArrayFlags", "-mapMultiple", "one:alpha", "-mapMultiple", "two:bravo"}
	var testFlags StringMap
	flag.Var(&testFlags, "mapMultiple", "Some description for this param.")
	flag.Parse()
	require.NotNil(t, testFlags)
	assert.Len(t, testFlags, 2)
	assert.Equal(t, "alpha", testFlags["one"])
	assert.Equal(t, "bravo", testFlags["two"])
}

func Example_mapStringFlags() {
	os.Args = []string{"TestStringArrayFlags", "-data", "one:alpha", "-data", "two:bravo"}
	var testFlags StringMap
	flag.Var(&testFlags, "data", "Some description for this param.")
	flag.Parse()
	fmt.Println(testFlags)
	// Output: map[one:alpha two:bravo]
}
