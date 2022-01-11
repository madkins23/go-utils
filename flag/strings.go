package flag

import (
	"regexp"
	"strings"
)

// StringArray defines a flag that can be invoked multiple times with values accumulated in an array.
// Comma-separated values may be combined in a single flag argument to be separated into the array.
// Extra space(s) around the commas are removed.
type StringArray []string

// String representation of the array of flag values.
func (i *StringArray) String() string {
	return "[" + strings.Join(*i, ",") + "]"
}

var commaSplitter = regexp.MustCompile("\\s*,\\s*")

// Set a value(s) into the array.
func (i *StringArray) Set(value string) error {
	*i = append(*i, commaSplitter.Split(value, -1)...)
	return nil
}
