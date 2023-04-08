package path

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

const (
	midFixPath = "goober"
	relFixPath = "snoofus.ext"
)

// HomePath tests.
// NOTE: Error conditions not testable since os/user isn't an interface.

func ExampleFixHomePath() {
	fixed, err := FixHomePath("~/big/bopper")
	if err == nil {
		fmt.Println(fixed)
	}
}

func TestFixHomePath_justHome(t *testing.T) {
	absPath, err := FixHomePath("~/")
	assert.NoError(t, err)
	assert.Equal(t, getUserHomeDir(t), absPath)
}

func TestFixHomePath_one(t *testing.T) {
	absPath, err := FixHomePath("~/" + relFixPath)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(getUserHomeDir(t), relFixPath), absPath)
}

func TestFixHomePath_multi(t *testing.T) {
	absPath, err := FixHomePath("~/"+midFixPath, relFixPath)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(getUserHomeDir(t), midFixPath, relFixPath), absPath)
}

// Get home directory via a different mechanism from that used by FixHomePath.
func getUserHomeDir(t *testing.T) string {
	homeDir, err := homedir.Dir()
	assert.NoError(t, err)
	return homeDir
}
