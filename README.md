# go-utils

Common Go packages shared by other Go projects.

These are not really production quality packages, nor are they maintained or supported.
They're simplistic mechanisms implemented to solve specific problems
in code that the author has been writing for his own use.

*If you need this functionality consider copying the code into your own project and modifying to fit your need.*

See the source for documentation.

## `path`

Path utilities.

* `path.HomePath()` prepends a relative file path with the user's home directory.
Works on linux and should work on Mac and Windows but untested by author.

## `test`

Test utilities.

* `test.CaptureStderr()` captures standard error over the execution
of a provided test function and returns the text so captured.
* `test.CaptureStdout()` captures standard output over the execution
of a provided test function and returns the text so captured.

## `text`

Text utilities.

* `text.Pluralize()` makes a word singular or plural based on the specified count.
Based on `github.com/gertd/go-pluralize` with hidden global pluralizer and
simplified calling convention.
* `text.JustAlphaNumeric()` filters non-alphanumeric characters out of a string.

## `typeutils`

Type utilities.

* `typeutils.Registry` provides a way to register types by name.
Normally Go doesn't keep type names at runtime, so it must be done by the application.
The `Registry` object provides a way to track this and to generate objects of a "named" type.
Created for use in Marshaling/Unmarshaling objects.
Uses reflection. Not thread-safe.
See test files for examples of JSON and YAML marshal/unmarshal.

* `typeutils.Registrar` provides a thread-safe `Registry`.
`Registry` methods are wrapped with a mutex object.
