package msg

//////////////////////////////////////////////////////////////////////////

var _ error = ConstError("poop")

// ConstError defines an Error type that can be a constant because it's really just a string.
type ConstError string

func (ce ConstError) Error() string {
	return string(ce)
}

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

// Is determines if the error is or contains the target error.
func (b *ErrBlocked) Is(target error) bool {
	if eb, ok := target.(*ErrBlocked); !ok {
		return false
	} else if eb.Name != "" {
		return eb.Name == b.Name
	} else {
		return true
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

// Is determines if the error is or contains the target error.
func (d *ErrDeprecated) Is(target error) bool {
	if ed, ok := target.(*ErrDeprecated); !ok {
		return false
	} else if ed.Name != "" {
		return ed.Name == d.Name
	} else {
		return true
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

// Is determines if the error is or contains the target error.
func (ni *ErrNotImplemented) Is(target error) bool {
	if eni, ok := target.(*ErrNotImplemented); !ok {
		return false
	} else if eni.Name != "" {
		return eni.Name == ni.Name
	} else {
		return true
	}
}
