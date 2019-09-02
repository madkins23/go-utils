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

## `text`

Text utilities.

* `text.Pluralize()` makes a word singular or plural based on the specified count.
Based on `github.com/gertd/go-pluralize` with hidden global pluralizer and
simplified calling convention.
