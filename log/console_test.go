package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestConsoleOrFile_Setup(t *testing.T) {
	cof := ConsoleOrFile{
		Console: true,
	}
	require.NoError(t, cof.Setup())
	cof = ConsoleOrFile{
		Console: false,
		LogFile: "/tmp/console-or-file.log",
		AsJSON:  true,
	}
	require.NoError(t, cof.Setup())
}

func ExampleConsole() {
	Console()
	Logger().Info().Msg("Easier to read")
	// Output:
}

func ExampleConsoleOrFile_formatExtra() {
	cof := ConsoleOrFile{
		Console:   true,
		UseStdout: true,
	}
	if err := cof.Setup(); err != nil {
		fmt.Printf("Log file creation error: %s\n", err)
		return
	}
	defer cof.CloseForDefer()

	// Get pointer to starting zerolog.ConsoleWriter().
	writer := cof.Writer()
	if writer == nil {
		fmt.Printf("No writer available\n")
		return
	}
	// Make a copy or the changes will persist.
	fixed := *writer
	// Make changes to the new ConsoleWriter.
	// Exclude timestamp so that example output will match.
	fixed.NoColor = true
	fixed.PartsExclude = []string{zerolog.TimestampFieldName}
	fixed.FormatExtra =
		// This function will add extra stuff to the end of each log line.
		// The data map passed in contains all the items to be logged by name.
		func(data map[string]interface{}, buffer *bytes.Buffer) error {
			// Delete timestamp from data so that example output will match.
			delete(data, zerolog.TimestampFieldName)
			if marshaled, err := json.MarshalIndent(data, "", "  "); err != nil {
				return fmt.Errorf("marshal data into json: %w", err)
			} else {
				buffer.WriteString("\n")
				buffer.Write(marshaled)
				buffer.WriteString("\n")
				return nil
			}
		}
	// Make a new log with the fixed ConsoleWriter.
	myLog := log.Output(fixed)
	myLog.Info().Str("alpha", "one").Int("bravo", 2).Msg("reversed")

	// Output:
	// INF reversed alpha=one bravo=2
	// {
	//   "alpha": "one",
	//   "bravo": 2,
	//   "level": "info",
	//   "message": "reversed"
	// }
}

func ExampleConsoleOrFile() {
	cof := ConsoleOrFile{
		Console: false,
		LogFile: "/tmp/console-or-file.log",
		AsJSON:  true,
	}
	if err := cof.Setup(); err != nil {
		fmt.Printf("Log file creation error: %s", err)
		return
	}
	defer cof.CloseForDefer()
	// Output:
}

func ExampleConsoleOrFile_withFlags() {
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	cof := ConsoleOrFile{}
	cof.AddFlagsToSet(flags, "/tmp/console-or-file.log")
	if err := flags.Parse([]string{"-console"}); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Printf("Error parsing command line flags: %s", err)
		}
		return
	}
	if err := cof.Setup(); err != nil {
		fmt.Printf("Log file creation error: %s", err)
		return
	}
	defer cof.CloseForDefer()
	// Output:
}
