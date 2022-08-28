package msg

//////////////////////////////////////////////////////////////////////////

const strBlocked = "blocked"
const strBlockedNamed = " is blocked"

// ErrBlocked is a custom error representing a blocked function or method.
type ErrBlocked struct {
	// Optional name of deprecated function or method.
	Name string
}

// Error implements the predefined error interface.
func (b *ErrBlocked) Error() string {
	if b.Name == "" {
		return strBlocked
	} else {
		return b.Name + strBlockedNamed
	}
}

//////////////////////////////////////////////////////////////////////////

const strDeprecated = "deprecated"
const strDeprecatedNamed = " is deprecated"

// ErrDeprecated is a custom error representing a deprecated function or method.
type ErrDeprecated struct {
	// Optional name of deprecated function or method.
	Name string
}

// Error implements the predefined error interface.
func (d *ErrDeprecated) Error() string {
	if d.Name == "" {
		return strDeprecated
	} else {
		return d.Name + strDeprecatedNamed
	}
}

//////////////////////////////////////////////////////////////////////////

const strNotImplementedYet = "not implemented yet"
const strNotImplementedNamed = " is not implemented yet"

// ErrNotImplemented is a custom error representing an unimplemented function or method.
type ErrNotImplemented struct {
	// Optional name of unimplemented function or method.
	Name string
}

// Error implements the predefined error interface.
func (ni *ErrNotImplemented) Error() string {
	if ni.Name == "" {
		return strNotImplementedYet
	} else {
		return ni.Name + strNotImplementedNamed
	}
}
