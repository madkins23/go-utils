// Package error defines some non-standard error mechanisms.
//
//  * Error with an error with an array of strings representing error details
//  * Error with an error with a map of strings representing error details
//  * Add NotYetImplemented() and ToBeImplemented(name) 'standard' error functions.
//
// The name of this package is somewhat unfortunate given that it conflicts with the error interface builtin.
// This package may be deprecated/renamed at some point if a better name ever presents itself.
// In the meantime just rename it in the use statement.
package error
