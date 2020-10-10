package reg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const aliasName = "test"

type example1 struct{}

type example2 struct{}

func TestAlias(t *testing.T) {
	alias := NewAlias(aliasName, nil)
	require.NotNil(t, alias)
	reg, ok := alias.Registry.(*registry)
	require.True(t, ok)
	require.NotNil(t, reg)
	assert.False(t, alias.aliased)
	assert.Len(t, reg.aliases, 0)
	assert.Len(t, reg.byName, 0)
	require.NoError(t, alias.Register(&example1{}))
	assert.True(t, alias.aliased)
	assert.Len(t, reg.aliases, 1)
	assert.Len(t, reg.byName, 1)
	// Since we can't redefine an alias (see registry_test.go)
	// doing it twice would generate an error here.
	require.NoError(t, alias.Register(&example2{}))
	assert.Len(t, reg.aliases, 1)
	assert.Len(t, reg.byName, 2)
}
