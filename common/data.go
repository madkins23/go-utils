package common

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	dockerDirectory = "/data"
	errNoDataDir    = "unable to locate data directory"
)

var (
	dataDirectory string
)

func DataDirectory(relativeDir string) (string, error) {
	if dataDirectory == "" {
		var err error

		if _, err = os.Stat(dockerDirectory); !os.IsNotExist(err) {
			dataDirectory = dockerDirectory
		} else {
			var workingDir string
			workingDir, err = os.Getwd()
			if err != nil {
				return "", errors.Wrap(err, errNoDataDir)
			}

			localDirectory := filepath.Join(workingDir, relativeDir, "data")
			if _, err = os.Stat(localDirectory); !os.IsNotExist(err) {
				dataDirectory = localDirectory
			}
		}

		if dataDirectory == "" {
			if err != nil {
				return "", errors.Wrap(err, errNoDataDir)
			} else {
				return "", errors.New(errNoDataDir)
			}
		}
	}

	return dataDirectory, nil
}
