package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatUSD(t *testing.T) {
	testFormatUSD(t, "$1,234,567.89", "1234567.89")
	testFormatUSD(t, "$234,567.89", "234567.89")
	testFormatUSD(t, "$34,567.89", "34567.89")
	testFormatUSD(t, "$4,567.89", "4567.89")
	testFormatUSD(t, "$567.89", "567.89")
	testFormatUSD(t, "$67.89", "67.89")
	testFormatUSD(t, "$7.89", "7.89")
	testFormatUSD(t, "$0.89", ".89")
	testFormatUSD(t, "$67.12", "67.127")
	testFormatUSD(t, "$67.12", "67.123")
	testFormatUSD(t, "$67.12", "67.12")
	testFormatUSD(t, "$67.10", "67.1")
	testFormatUSD(t, "$67.00", "67")
	testFormatUSD(t, "$7.00", "7")
	testFormatUSD(t, "$0.00", "")
	testFormatUSD(t, "$0.00", ".")
}

func testFormatUSD(t *testing.T, expected, number string) {
	formatted, err := FormatUSD(number)
	require.NoError(t, err)
	assert.Equal(t, expected, formatted)
	formatted, err = FormatUSD("-" + number)
	require.NoError(t, err)
	assert.Equal(t, "-"+expected, formatted)
}
