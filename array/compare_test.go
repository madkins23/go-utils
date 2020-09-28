package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringElementsMatch(t *testing.T) {
	one := []string{"alpha", "bravo", "charlie"}
	two := []string{"alpha", "bravo", "charlie"}
	three := []string{"alpha", "charlie", "bravo"}
	four := []string{"charlie", "bravo", "alpha"}
	five := []string{"alpha", "bravo", "charlie", "delta"}
	six := []string{"alpha", "bravo"}
	seven := []string{"alpha", "bravo", "delta"}
	assert.True(t, StringElementsMatch(nil, nil))
	assert.True(t, StringElementsMatch([]string{}, nil))
	assert.True(t, StringElementsMatch([]string{}, []string{}))
	assert.True(t, StringElementsMatch(one, two))
	assert.True(t, StringElementsMatch(one, three))
	assert.True(t, StringElementsMatch(one, four))
	assert.False(t, StringElementsMatch(one, five))
	assert.False(t, StringElementsMatch(one, six))
	assert.False(t, StringElementsMatch(one, seven))
}
