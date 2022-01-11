package flag

import (
	"errors"
	"regexp"
	"strings"
)

// StringMap defines a flag that can be invoked multiple times with values accumulated in a map.
// Each flag value is broken at its embedded colon to create a key/value pair.
type StringMap map[string]string

// String representation of the array of flag values.
func (i *StringMap) String() string {
	pairs := make([]string, 0, len(*i))
	for k, v := range *i {
		pairs = append(pairs, k+":"+v)
	}
	return "[" + strings.Join(pairs, ",") + "]"
}

var colonSplitter = regexp.MustCompile("\\s*:\\s*")
var errNoColon = errors.New("no colon")

// Set a key/value item into the array.
func (i *StringMap) Set(value string) error {
	stuff := colonSplitter.Split(value, 2)
	if len(stuff) < 2 {
		return errNoColon
	}
	if *i == nil {
		*i = make(StringMap)
	}
	(*i)[stuff[0]] = stuff[1]
	return nil
}
