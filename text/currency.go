package text

import (
	"fmt"
)

const dollar = '$'

// FormatUSD takes a string representation of a number,
// decorates it with a preceding dollar sign, inserts commas,
// and adds trailing pennies as required.
func FormatUSD(number string) (string, error) {
	if separated, err := AddNumericSeparators(number, ','); err != nil {
		return "", fmt.Errorf("add comma separators: %w", err)
	} else {
		return "$" + separated, nil
	}
}
