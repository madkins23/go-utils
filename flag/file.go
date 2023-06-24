package flag

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/madkins23/go-utils/path"
)

// LoadSettings overrides default values in the specified flag.FlagSet with
// values taken from a settings file(s) specified as '@<path>' on the command line.
// All settings file specifications are removed from os.Args by this function.
// Each settings file is unmarshaled into a map[string]string:
// JSON files contain legitimate JSON, CFG files contain 'key = value' lines.
// Each key in the map is looked up in the specified flagSet and its default value
// and current value replaced by the value associated with that key from the settings file.
// After this the flagSet can be parsed so that the flag values from the command
// line override the defaults from the settings file as needed.
// Order of settings files is important, later ones override earlier ones.
// Flags remaining in os.Args can then be processed in the normal manner.
func LoadSettings(flagSet *flag.FlagSet) error {
	fixedArgs := make([]string, 0, len(os.Args))
	for _, arg := range os.Args {
		if !strings.HasPrefix(arg, "@") {
			fixedArgs = append(fixedArgs, arg)
		} else if err := loadSettingsFile(flagSet, arg[1:]); err != nil {
			return fmt.Errorf("load settings file: %w", err)
		}
	}
	os.Args = fixedArgs

	return nil
}

var errEmptyPath = errors.New("empty path")

func loadSettingsFile(flagSet *flag.FlagSet, file string) error {
	if file == "" {
		return errEmptyPath
	}

	// TODO: Handle different types of files (not just JSON)
	var settings map[string]string
	if fixed, err := path.FixHomePath(file); err != nil {
		return fmt.Errorf("fix path '%s': %w", file, err)
	} else if bytes, err := os.ReadFile(fixed); err != nil {
		return fmt.Errorf("read file '%s': %w", fixed, err)
	} else {
		var ext = filepath.Ext(file)
		switch strings.ToLower(ext) {
		case ".json":
			if err = json.Unmarshal(bytes, &settings); err != nil {
				return fmt.Errorf("unmarshal JSON settings: %w", err)
			}
		case ".cfg":
			if settings, err = configUnmarshal(string(bytes)); err != nil {
				return fmt.Errorf("unmarshal CFG settings: %w", err)
			}
		default:
			return fmt.Errorf("unknown settings file extension '%s'", ext)
		}
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

	return nil
}

var errBadRegex = errors.New("unable to compile regex")

func configUnmarshal(text string) (map[string]string, error) {
	if re, err := regexp.Compile("^\\s*(\\S+)\\s*=\\s*(.*)\\s*$"); err != nil {
		return nil, errBadRegex
	} else {
		settings := make(map[string]string)
		for _, line := range strings.Split(text, "\n") {
			if subMatches := re.FindStringSubmatch(line); subMatches != nil {
				settings[subMatches[1]] = subMatches[2]
			} else if strings.Trim(line, " \t") != "" {
				return nil, fmt.Errorf("unknown config line '%s'", line)
			}
		}
		return settings, nil
	}
}
