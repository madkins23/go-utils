package path

import (
	"os/user"
	"path/filepath"

	"golang.org/x/xerrors"
)

// Return a path constructed from the specified relative path and the user's home directory.
func homePath(relPath string) (string, error) {
	// NOTE: Error conditions not easily testable since os/user isn't an interface.
	usr, err := user.Current()
	if err != nil || usr == nil {
		return "", xerrors.Errorf("No user data getting home path: %w", err)
	}
	if usr.HomeDir == "" {
		return "", xerrors.Errorf("No home directory for user: %w", err)
	}

	return filepath.Join(usr.HomeDir, relPath), nil
}
