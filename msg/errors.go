package msg

//////////////////////////////////////////////////////////////////////////

const strBlocked = "blocked"
const strBlockedNamed = " is blocked"

// ErrBlocked is a custom error representing a blocked function or method.
type ErrBlocked struct {
	// Optional name of deprecated function or method.
	name string
}

// Error implements the predefined error interface.
func (b *ErrBlocked) Error() string {
	if b.name == "" {
		return strBlocked
	} else {
		return b.name + strBlockedNamed
	}
}

//////////////////////////////////////////////////////////////////////////

const strDeprecated = "deprecated"
const strDeprecatedNamed = " is deprecated"

// ErrDeprecated is a custom error representing a deprecated function or method.
type ErrDeprecated struct {
	// Optional name of deprecated function or method.
	name string
}

// Error implements the predefined error interface.
func (d *ErrDeprecated) Error() string {
	if d.name == "" {
		return strDeprecated
	} else {
		return d.name + strDeprecatedNamed
	}
}

//////////////////////////////////////////////////////////////////////////

const strNotImplementedYet = "not implemented yet"
const strNotImplementedNamed = " is not implemented yet"

// ErrNotImplemented is a custom error representing an unimplemented function or method.
type ErrNotImplemented struct {
	// Optional name of unimplemented function or method.
	name string
}

// Error implements the predefined error interface.
func (ni *ErrNotImplemented) Error() string {
	if ni.name == "" {
		return strNotImplementedYet
	} else {
		return ni.name + strNotImplementedNamed
	}
}
