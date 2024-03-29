# go-utils

Common Go packages shared by other Go projects.

See the [source](https://github.com/madkins23/go-utils)
or [godoc](https://godoc.org/github.com/madkins23/go-utils) for more detailed documentation.

[![Go Report Card](https://goreportcard.com/badge/github.com/madkins23/go-utils)](https://goreportcard.com/report/github.com/madkins23/go-utils)
![GitHub](https://img.shields.io/github/license/madkins23/go-utils)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/madkins23/go-utils)
[![Go Reference](https://pkg.go.dev/badge/github.com/madkins23/go-utils.svg)](https://pkg.go.dev/github.com/madkins23/go-utils)

## `app`

Application utilities.

* `app.Terminator` provides support for "graceful" shutdown.
* `app.HandleSignals` invokes a `SignalHandler` when one of the specified
  OS signals is invoked.
* `app.HandleTerminalSignals` invokes a `SignalHandler` when an OS signal
  that normally causes application termination is invoked.

## `array`

Array utilities.

* `array.StringElementsMatch()` compares two arrays to see if they match irrespective of order.

## `check`

* `check.IsZero()` and `check.ErrorIfZero()` use reflection to check if an entity
  has the zero value for its type.

## `cycle`

Periodic code execution.

* `cycle.Periodic` type provides a mechanism for cyclically executing code.
* `Periodic.Ticker` executes code at specified intervals.

## `csv`

Wrapper around `encoding/csv.Reader` to provide:

* header-aware processing and
* map object result (header -> value).

## `error`

**Deprecated**

Badly named package has been replaced by features in the `msg` package
and the new-ish Go package feature `errors.Join()`.

## `flag`

Extend argument parsing behavior.

* `flag.LoadSettings()` parses flag configuration files.
* `flag.StringMap` defines a flag that can be invoked multiple times with values accumulated in a map.
* `flag.StringArray` defines a flag that can be invoked multiple times with values accumulated in an array.

## `log`

Logging utilities using [zerolog](https://github.com/rs/zerolog).

* `log.Logger()` returns the default zerolog logger.
* `log.LocalLogger` is a logging mixin for embedding in other structs.
* `log.Console` configures the default zerolog logger for readable format
  instead of the default JSON record output.

## `msg`

Error utilities.
Putting them in a package named `error` or `errors` makes life difficult
so the package is named `msg`.

* `msg.ConstError` is a string type that implements the `error` interface.
* several general-purpose error messages implemented as `struct` items:
  * `msg.ErrBlocked`
  * `msg.ErrDeprecated`
  * `msg.ErrNotImplemented`

## `path`

Path utilities.

* `path.FixHomePath()` looks for a home directory string and replaces it with 
  the absolute path to the home directory.
  Works on linux and should work on Mac and Windows but untested by author.
* `path.HomePath()` prepends a relative file path with the user's home directory.
  Works on linux and should work on Mac and Windows but untested by author.
  _This function has been deprecated in favor of `path.FixHomePath()`._

## `server`

Utilities for interacting with Golang HTTP servers.

* `server.Interrupt` sends a `SIGINT` signal to the current process.
* `server.IsReady` checks for immediate service of specified URL.
* `server.WaitFor` waits for service to become available for specified timeout.

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
* `text.AddNumericSeparators` adds separators (commas or whatever)
  every three digits in the integer part of a number.
* `text.FormatUSD()` invokes `text.AddNumericSeparators` and adds a dollar sign prefix.
* `text.TableDef` defines a simple table structure that can be used to generate:
  * a format string for the header row of the table,
  * a divider string for use below the header or in between table sections, and
  * a format string for the data rows of the table.

## `typeutils`

**Deprecated**

Original location of type registration mechanism.
This code has since been removed to the `reg` package in the
[`go-type`](https://github.com/madkins23/go-type) project.
The `typeutils` package will be removed entirely in any future `V2` version.
