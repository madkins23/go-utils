package text

import "github.com/gertd/go-pluralize"

var (
	// Persistent go-pluralize object.
	plural = pluralize.NewClient()
)

// Pluralize returns singular or plural form of the specified word based on the count.
// Count of 1 is singular, otherwise the plural form is returned.
// Uses the gertd/go-pluralize/Pluralize() function, just simplifies usage a tad.
// May only work if the word is singular.
func Pluralize(word string, count int) string {
	return plural.Pluralize(word, count, false)
}
