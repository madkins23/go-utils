package text

import (
	"log"
	"regexp"
)

const (
	patternAlphaNumeric = "[^a-zA-Z0-9]+"
)

var (
	regexAphaNumeric *regexp.Regexp
)

func init() {
	// Initialize for JustAlphaNumeric()
	var err error
	regexAphaNumeric, err = regexp.Compile(patternAlphaNumeric)
	if err != nil {
		// Should never happen and there's nothing we can do to fix it.
		log.Fatalf("*** Unable to compile alphanumeric regex: %v", err)
	}
}

// justAlphaNumeric removes all non-alphanumeric characters from a string.
func JustAlphaNumeric(text string) string {
	return regexAphaNumeric.ReplaceAllString(text, "")
}
