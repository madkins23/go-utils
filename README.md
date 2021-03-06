# go-utils

Common Go packages shared by other Go projects.

You are more than welcome to use this package as is but these are
utility packages constructed by the author for use in personal projects.
The author makes occasional changes and attempts to follow proper versioning and release protocols,
however this code should not be considered production quality or maintained.

*Consider copying the code into your own project and modifying to fit your need.*

See the [source](https://github.com/madkins23/go-utils)
or [godoc](https://godoc.org/github.com/madkins23/go-utils) for documentation.

## `array`

Array utilities.

* `array.StringElementsMatch()` compares two arrays to see if they match irrespective of order.

## `cycle`

Periodic code execution.

* `cycle.Periodic` type provides a mechanism for cyclically executing code.
* `Periodic.Ticker` executes code at specified intervals.

## `log`

Logging utilities using [zerolog](https://github.com/rs/zerolog).

* `log.Logger()` returns the default zerolog logger.
* `log.LocalLogger` is a logging mixin for embedding in other structs.
* `log.Console` configures the default zerolog logger for readable format
  instead of the default JSON record output.

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

* `text.JustAlphaNumeric()` filters non-alphanumeric characters out of a string.
* `text.Pluralize()` makes a word singular or plural based on the specified count.
  Calls through to `github.com/gertd/go-pluralize` with hidden global pluralizer
  and simplified calling convention.

## `typeutils`

Original location of type registration mechanism.
This code has since been removed to the `reg` package in the
[`go-type`](https://github.com/madkins23/go-type) project.
The `typeutils` package will be removed entirely in any future `V2` version.
