package flag

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSettings(t *testing.T) {
	var (
		text  string
		whole int
		float float64
	)
	flagSet := flag.NewFlagSet("example", flag.ContinueOnError)
	flagSet.StringVar(&text, "text", "Lorem Ipsum", "Text String")
	flagSet.IntVar(&whole, "whole", 13, "Integer")
	flagSet.Float64Var(&float, "float", 1.61803, "Floating point")

	// Original defaults when no settings file or flags:
	require.NoError(t, flagSet.Parse([]string{}))
	assert.Equal(t, "Lorem Ipsum", text)
	assert.Equal(t, 13, whole)
	assert.Equal(t, 1.61803, float)

	// Add the settings file to the command line arguments.
	os.Args = []string{"path", "@testdata/settings.json"}
	require.NoError(t, LoadSettings(flagSet))

	// Settings override defaults:
	require.NoError(t, flagSet.Parse(os.Args[1:]))
	assert.Equal(t, "Don't Look!", text)
	assert.Equal(t, 17, whole)
	assert.Equal(t, 2.71828, float)

	// Use the settings file And the command line arguments.
	os.Args = []string{
		"path", "@testdata/settings.json",
		"-text", "Read Me!",
		"-whole", "23",
		"-float", "3.14159",
	}
	require.NoError(t, LoadSettings(flagSet))

	// Flags override settings and/or defaults:
	require.NoError(t, flagSet.Parse(os.Args[1:]))
	assert.Equal(t, "Read Me!", text)
	assert.Equal(t, 23, whole)
	assert.Equal(t, 3.14159, float)
}

func TestLoadSettings_badPath(t *testing.T) {
	flagSet := flag.NewFlagSet("example", flag.ContinueOnError)
	os.Args = []string{"path", "@testdata/noSuchFile.json"}
	assert.Error(t, LoadSettings(flagSet))
}

func TestLoadSettings_badJSON(t *testing.T) {
	flagSet := flag.NewFlagSet("example", flag.ContinueOnError)
	os.Args = []string{"path", "@testdata/badFile.json"}
	assert.Error(t, LoadSettings(flagSet))
}

func TestLoadSettings_badSetting(t *testing.T) {
	flagSet := flag.NewFlagSet("example", flag.ContinueOnError)
	os.Args = []string{"path", "@testdata/settings.json"}
	assert.Error(t, LoadSettings(flagSet))
}

func TestLoadSettings_badValue(t *testing.T) {
	flagSet := flag.NewFlagSet("example", flag.ContinueOnError)
	os.Args = []string{"path", "@testdata/badValue.json"}
	assert.Error(t, LoadSettings(flagSet))
}

////////////////////////////////////////////////////////////////////////////////

func ExampleLoadSettings() {
	var (
		err   error
		text  string
		whole int
		float float64
	)
	flagSet := flag.NewFlagSet("example", flag.ContinueOnError)
	flagSet.StringVar(&text, "text", "Lorem Ipsum", "Text String")
	flagSet.IntVar(&whole, "whole", 13, "Integer")
	flagSet.Float64Var(&float, "float", 1.61803, "Floating point")

	os.Args = []string{
		"path", "@testdata/settings.json",
		"-text", "Read Me!",
		"-whole", "23",
		"-float", "3.14159",
	}
	err = LoadSettings(flagSet)
	if err != nil {
		fmt.Printf("Error loading flag settings: %s\n", err)
		return
	}

	// Flags override settings and/or defaults:
	if err = flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error parsing flags: %s\n", err)
		return
	}
	fmt.Printf("%s %d %7.5f\n", text, whole, float)

	// Output: Read Me! 23 3.14159
}
