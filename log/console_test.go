package log

import (
	"errors"
	"flag"
	"fmt"
)

func ExampleConsole() {
	Console()
	Logger().Info().Msg("Easier to read")
	// Output:
}

func ExampleConsoleOrFile() {
	flags := flag.NewFlagSet("gardepro", flag.ContinueOnError)
	console := flags.Bool("console", false, "Direct log to console")
	logFile := flags.String("log", "/tmp/console-or-file.log", "Path to log file")
	logJSON := flags.Bool("logJSON", false, "Log output as JSON")
	if err := flags.Parse([]string{"-console"}); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Printf("Error parsing command line flags: %s", err)
		}
		return
	}

	cof := ConsoleOrFile{
		Console: *console,
		LogFile: *logFile,
		AsJSON:  *logJSON,
	}
	if err := cof.Setup(); err != nil {
		fmt.Printf("Log file creation error: %s", err)
		return
	}
	defer cof.CloseForDefer()
	// Output:
}
