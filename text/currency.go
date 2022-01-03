package text

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var dollar = '$'
var period = '.'
var multipleDecimalPoints = errors.New("multiple decimal points")

// FormatUSD takes a string representation of a number,
// decorates it with a preceding dollar sign, and inserts commas as required.
//
// This is not intended to be a general purpose solution,
// it is a stopgap until golang formalizes decimal and currency libraries.
// In particular, shopspring/decimal.Decimal doesn't have formatting with embedded commas.
// Formatting for decimal data (if/when added to Golang) may be in the text package.
//
// The current implementation isn't terribly efficient.
// There may be a state-machine solution that checks through the number string in reverse.
func FormatUSD(number string) (string, error) {
	decimalFound := false
	intPart := strings.Builder{}
	intPart.Grow(len(number))
	decPart := strings.Builder{}
	decPart.Grow(len(number))

	for _, c := range number {
		if c == period {
			if decimalFound {
				// Can't have two dots
				return "", multipleDecimalPoints
			}
			decimalFound = true
		} else if !unicode.IsDigit(c) {
			return "", fmt.Errorf("unknown numeric character '%c'", c)
		} else if decimalFound {
			decPart.WriteRune(c)
		} else {
			intPart.WriteRune(c)
		}
	}

	result := strings.Builder{}
	result.Grow(2 * len(number))
	result.WriteRune(dollar)

	var whole string
	if intPart.Len() > 0 {
		whole = intPart.String()
	} else {
		whole = "0"
	}

	comma := false
	if xtra := len(whole) % 3; xtra != 0 {
		result.WriteString(whole[0:xtra])
		whole = whole[xtra:]
		comma = true
	}

	for len(whole) > 0 {
		if comma {
			result.WriteRune(',')
		} else {
			comma = true
		}
		result.WriteString(whole[0:3])
		whole = whole[3:]
	}

	var dec string
	result.WriteRune(period)
	if decimalFound {
		dec = decPart.String()
		if len(dec) > 2 {
			dec = dec[0:2]
		}
		result.Write([]byte(dec))
	}
	for i := len(dec); i < 2; i++ {
		result.WriteRune('0')
	}

	return result.String(), nil
}
