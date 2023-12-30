package text

import (
	"math"
	"strings"

	"github.com/madkins23/go-utils/msg"
)

const (
	minus  = '-'
	period = '.'
	zero   = '0'
)

var multipleDecimalPoints msg.ConstError = "multiple decimal points"

// AddNumericSeparators takes a string representation of a number,
// inserts the specified separator character and adds trailing decimal data as required.
func AddNumericSeparators(number string, separator rune) (string, error) {
	return addNumericSeparators(number, separator, 0, -1)
}

func addNumericSeparators(number string, separator rune, decimalRequired int, decimalLimit int) (string, error) {
	prefix := ""
	if len(number) > 0 && number[0] == minus {
		prefix = "-"
		number = number[1:]
	}

	parts := strings.Split(number, string(period))
	if len(parts) > 2 {
		return number, multipleDecimalPoints
	}

	var decPart string
	intPart := parts[0]
	if len(parts) > 1 {
		decPart = parts[1]
		if decimalLimit >= 0 && len(decPart) > decimalLimit {
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

	if decimalRequired > 0 || (len(parts) > 1 && len(decPart) > 0) {
		result.WriteRune(period)

		result.WriteString(decPart)
		if decimalRequired > 0 {
			for i := len(decPart); i < decimalRequired; i++ {
				result.WriteRune(zero)
			}
		}
	}

	return prefix + result.String(), nil
}
