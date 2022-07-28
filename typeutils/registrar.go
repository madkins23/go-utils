package typeutils

import "sync"

// NewRegistrar creates a new Registrar object of the default internal type.
// Registries created via this function are mutex locked for concurrent access.
//
// Deprecated: This functionality has been rewritten in madkins23/go-type
func NewRegistrar() Registry {
	return &registrar{
		Registry: NewRegistry(),
	}
}

//////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////

// Default Registrar implementation.
type registrar struct {
	Registry
	lock sync.Mutex
}

// Alias creates an alias to be used to shorten names.
// Use an empty string to remove a previous alias.
//
// Deprecated: This functionality has been rewritten in madkins23/go-type
func (reg *registrar) Alias(alias string, example interface{}) error {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return reg.Registry.Alias(alias, example)
}

// Register a type by providing an example object.
//
// Deprecated: This functionality has been rewritten in madkins23/go-type
func (reg *registrar) Register(example interface{}) error {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return reg.Registry.Register(example)
}

// Make creates a new instance of the example object with the specified name.
// The new instance will be created with fields filled with zero values.
//
// Deprecated: This functionality has been rewritten in madkins23/go-type
func (reg *registrar) Make(name string) (interface{}, error) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return reg.Registry.Make(name)
}

// NameFor returns a name for the specified object.
//
// Deprecated: This functionality has been rewritten in madkins23/go-type
func (reg *registrar) NameFor(item interface{}) (string, error) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	return reg.Registry.NameFor(item)
}
