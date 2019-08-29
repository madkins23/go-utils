package path

import (
	"os/user"
	"path/filepath"

	"golang.org/x/xerrors"
)

var (
	ErrNoUserData      = xerrors.New("no user data")
	ErrNoHomeDirectory = xerrors.New("no home directory for user")
)

// Return a path constructed from the specified relative path and the user's home directory.
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
		return "", xerrors.Errorf("getting absolute path for ~/%s: %w", relPath, err)
	}

	return filepath.Join(append([]string{usr.HomeDir}, relPath...)...), nil
}
