package text

import (
	"math"
	"strings"

	"github.com/madkins23/go-utils/msg"
)

const (
	period = '.'
	zero   = '0'
)

var multipleDecimalPoints msg.ConstError = "multiple decimal points"

// AddNumericSeparators takes a string representation of a number,
// inserts the specified separator character and adds trailing decimal data as required.
func AddNumericSeparators(number string, separator rune) (string, error) {
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
	result.Grow(len(intPart) + len(intPart)/3 + 1 + int(math.Min(float64(len(decPart)), 2)))

	if intPart == "" {
		result.WriteRune(zero)
	} else {
		needSeparator := false
		if extra := len(intPart) % 3; extra != 0 {
			result.WriteString(intPart[0:extra])
			intPart = intPart[extra:]
			needSeparator = true
		}

		for len(intPart) > 0 {
			if needSeparator {
				result.WriteRune(separator)
			} else {
				needSeparator = true
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
