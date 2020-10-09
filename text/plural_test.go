package text

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamplePluralize() {
	fmt.Println(Pluralize("rat", 0))
	fmt.Println(Pluralize("rat", 1))
	fmt.Println(Pluralize("rat", 2))
	fmt.Println(Pluralize("ox", 0))
	fmt.Println(Pluralize("ox", 1))
	fmt.Println(Pluralize("ox", 2))
	// Output: rats
	// rat
	// rats
	// oxen
	// ox
	// oxen
}

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
