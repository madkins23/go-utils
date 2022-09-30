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
