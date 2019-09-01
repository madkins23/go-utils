package path

import (
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

const (
	midPath = "goober"
	relPath = "snoofus.ext"
)

// HomePath tests.
// NOTE: Error conditions not testable since os/user isn't an interface.

func TestHomePath_justHome(t *testing.T) {
	absPath, err := HomePath()
	assert.NoError(t, err)
	assert.Equal(t, getHomeDir(t), absPath)
}

func TestHomePath_one(t *testing.T) {
	absPath, err := HomePath(relPath)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(getHomeDir(t), relPath), absPath)
}

func TestHomePath_multi(t *testing.T) {
	absPath, err := HomePath(midPath, relPath)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(getHomeDir(t), midPath, relPath), absPath)
}

// Get home directory via a different mechanism from that used by HomePath.
func getHomeDir(t *testing.T) string {
	homeDir, err := homedir.Dir()
	assert.NoError(t, err)
	return homeDir
}
