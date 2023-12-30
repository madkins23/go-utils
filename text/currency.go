package text

import (
	"fmt"
)

// FormatUSD takes a string representation of a number,
// decorates it with a preceding dollar sign, inserts commas,
// and adds trailing pennies as required.
func FormatUSD(number string) (string, error) {
	prefix := "$"
	if len(number) > 0 && number[0] == minus {
		prefix = "-$"
		number = number[1:]
	}
	if separated, err := AddNumericSeparators(number, ','); err != nil {
		return "", fmt.Errorf("add comma separators: %w", err)
	} else {
		return prefix + separated, nil
	}
}
