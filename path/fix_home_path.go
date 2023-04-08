package path

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var errNoPath = errors.New("no path provided")

// FixHomePath returns the provided path adjusted to replace the home directory tilde if there is one.
// The path may be provided as a single string or a list of elements that will be joined.
func FixHomePath(path ...string) (string, error) {
	if len(path) < 1 {
		return "", errNoPath
	}

	if strings.HasPrefix(path[0], "~/") {
		if homeDir, err := os.UserHomeDir(); err != nil {
			return "", fmt.Errorf("no user homeDir directory: %w", err)
		} else {
			// Prepend home directory and fixed first path element to rest of path elements.
			path = append([]string{homeDir, path[0][2:]}, path[1:]...)
		}
	}

	return filepath.Join(path...), nil
}
