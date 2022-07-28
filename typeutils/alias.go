package typeutils

import "fmt"

type Alias struct {
	alias    string
	aliased  bool
	registry Registry
}

// NewAlias creates a new Alias.
//
// Deprecated: This functionality has been rewritten in madkins23/go-type
func NewAlias(registry Registry, alias string) *Alias {
	return &Alias{
		alias:    alias,
		registry: registry,
	}

}

// Register an alias using the specified example.
//
// Deprecated: This functionality has been rewritten in madkins23/go-type
func (a *Alias) Register(example interface{}) error {
	if !a.aliased {
		fmt.Println("nar", a.alias)
		if err := a.registry.Alias(a.alias, example); err != nil {
			return fmt.Errorf("register alias: %w", err)
		}
		a.aliased = true
	}

	if err := a.registry.Register(example); err != nil {
		return fmt.Errorf("register example: %w", err)
	}

	return nil
}
