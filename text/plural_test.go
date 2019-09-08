package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluralize_worms(t *testing.T) {
	assert.Equal(t, "worms", Pluralize("worm", 0))
	assert.Equal(t, "worm", Pluralize("worm", 1))
	assert.Equal(t, "worms", Pluralize("worm", 2))

	assert.Equal(t, "worms", Pluralize("worms", 0))
	assert.Equal(t, "worm", Pluralize("worms", 1))
	assert.Equal(t, "worms", Pluralize("worms", 2))

}

func TestPluralize_oxen(t *testing.T) {
	assert.Equal(t, "oxen", Pluralize("ox", 0))
	assert.Equal(t, "ox", Pluralize("ox", 1))
	assert.Equal(t, "oxen", Pluralize("ox", 2))

	assert.Equal(t, "oxen", Pluralize("oxen", 0))
	assert.Equal(t, "ox", Pluralize("oxen", 1))
	assert.Equal(t, "oxen", Pluralize("oxen", 2))

}
