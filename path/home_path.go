package path

import (
	"errors"
	"fmt"
	"os/user"
	"path/filepath"
)

var (
	// ErrNoUserData is returned when no user data can be acquired via user.Current().
	// This should never happen.
	ErrNoUserData = errors.New("no user data")

	// ErrNoHomeDirectory is returned when there is no home directory for the user.
	// This seems unlikely and will probably never happen.
	ErrNoHomeDirectory = errors.New("no home directory for user")

	// Since user.Current() is package level (i.e. there is no interface)
	// it's pretty much impossible to test for these conditions.
)

// HomePath returns a path constructed from the specified relative path and the user's home directory.
// Consider using https://github.com/mitchellh/go-homedir as an alternative.
//
// Deprecated: Prefer FixHomePath() in this package.
func HomePath(relPath ...string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		// Use this error.
	} else if usr == nil {
		err = ErrNoUserData
	} else if usr.HomeDir == "" {
		err = ErrNoHomeDirectory
	}
	if err != nil {
		return "", fmt.Errorf("getting absolute path for ~/%s: %w", relPath, err)
	}

	return filepath.Join(append([]string{usr.HomeDir}, relPath...)...), nil
}
