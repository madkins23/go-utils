package path

import (
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

const (
	relPath = "goober/snoofus.ext"
)

func TestHomePath(t *testing.T) {
	absPath, err := homePath(relPath)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(getHomeDir(t), relPath), absPath)
}

// Get home directory via a different mechanism from that used by homePath.
func getHomeDir(t *testing.T) string {
	homeDir, err := homedir.Dir()
	assert.NoError(t, err)
	return homeDir
}
