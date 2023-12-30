package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddNumericSeparators(t *testing.T) {
	testAddNumericSeparators(t, "1,234,567.89", "1234567.89")
	testAddNumericSeparators(t, "234,567.89", "234567.89")
	testAddNumericSeparators(t, "34,567.89", "34567.89")
	testAddNumericSeparators(t, "4,567.89", "4567.89")
	testAddNumericSeparators(t, "567.89", "567.89")
	testAddNumericSeparators(t, "67.89", "67.89")
	testAddNumericSeparators(t, "7.89", "7.89")
	testAddNumericSeparators(t, "0.89", ".89")
	testAddNumericSeparators(t, "67.123", "67.123")
	testAddNumericSeparators(t, "67.12", "67.12")
	testAddNumericSeparators(t, "67.1", "67.1")
	testAddNumericSeparators(t, "67", "67")
	testAddNumericSeparators(t, "7", "7")
	testAddNumericSeparators(t, "0", "")
	testAddNumericSeparators(t, "0", ".")
}

func testAddNumericSeparators(t *testing.T, expected, number string) {
	formatted, err := AddNumericSeparators(number, ',')
	require.NoError(t, err)
	assert.Equal(t, expected, formatted)
	formatted, err = AddNumericSeparators("-"+number, ',')
	require.NoError(t, err)
	assert.Equal(t, "-"+expected, formatted)
}
