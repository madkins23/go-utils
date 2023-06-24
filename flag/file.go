package flag

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/madkins23/go-utils/path"
)

const (
	SettingsFlag  = "settings"
	SettingsUsage = "Optional JSON file with flag defaults"
)

// LoadSettings overrides default values in the specified flag.FlagSet with
// values taken from a settings file specified as '@<path>' on the command line.
// The settings file specification is removed from os.Args by this function.
// If a settings file is found and is the path of a JSON file containing a map[string]string hash
// then each key is looked up in the specified flagSet and its default value
// and current value replaced by the value associated with that key from the settings file.
// After this the flagSet can be parsed so that the flag values from the command
// line override the defaults from the settings file as needed.
func LoadSettings(flagSet *flag.FlagSet) error {
	var file string
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "@") {
			file = arg[1:]
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
		}
	}

	if file != "" {
		var settings map[string]string
		if fixed, err := path.FixHomePath(file); err != nil {
			return fmt.Errorf("fix path '%s': %w", file, err)
		} else if bytes, err := os.ReadFile(fixed); err != nil {
			return fmt.Errorf("read file '%s': %w", fixed, err)
		} else if err = json.Unmarshal(bytes, &settings); err != nil {
			return fmt.Errorf("unmarshal settings: %w", err)
		}

		errs := make([]error, 0)
		for key, value := range settings {
			if flg := flagSet.Lookup(key); flg == nil {
				errs = append(errs, fmt.Errorf("setting referenced non-existent flag '%s'", key))
			} else {
				flg.DefValue = value
				if err := flg.Value.Set(value); err != nil {
					errs = append(errs, fmt.Errorf("set value of flag to '%s': %w", value, err))
				}
			}
		}
		if len(errs) > 0 {
			return fmt.Errorf("change defaults: %w", errors.Join(errs...))
		}
	}

	return nil
}
