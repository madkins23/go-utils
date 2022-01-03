package text

import (
	"errors"
	"math"
	"strings"
)

const dollar = '$'
const period = '.'
const zero = '0'

var multipleDecimalPoints = errors.New("multiple decimal points")

// FormatUSD takes a string representation of a number,
// decorates it with a preceding dollar sign, inserts commas,
// and adds trailing pennies as required.
//
// This is not intended to be a general purpose solution,
// it is a stopgap until golang formalizes decimal and currency libraries.
// In particular, shopspring/decimal.Decimal doesn't have formatting with embedded commas.
// Formatting for decimal data (if/when added to Golang) may be in the text package.
func FormatUSD(number string) (string, error) {
	parts := strings.Split(number, string(period))
	if len(parts) > 2 {
		return number, multipleDecimalPoints
	}

	var decPart string
	intPart := parts[0]
	if len(parts) > 1 {
		decPart = parts[1]
		if len(decPart) > 2 {
			decPart = decPart[0:2]
		}
	}

	result := strings.Builder{}
	result.Grow(1 + len(intPart) + len(intPart)/3 + 1 + int(math.Min(float64(len(decPart)), 2)))
	result.WriteRune(dollar)

	if intPart == "" {
		result.WriteRune(zero)
	} else {
		comma := false
		if extra := len(intPart) % 3; extra != 0 {
			result.WriteString(intPart[0:extra])
			intPart = intPart[extra:]
			comma = true
		}

		for len(intPart) > 0 {
			if comma {
				result.WriteRune(',')
			} else {
				comma = true
			}
			result.WriteString(intPart[0:3])
			intPart = intPart[3:]
		}
	}

	result.WriteRune(period)

	result.WriteString(decPart)
	for i := len(decPart); i < 2; i++ {
		result.WriteRune(zero)
	}

	return result.String(), nil
}
