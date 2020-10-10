package reg

import "fmt"

// Type Alias provides a package-specific alias and Registry combination.
// This simplifies regisration of types in a package with a common alias.
type Alias struct {
	Registry
	alias   string
	aliased bool
}

// NewAlias returns a package-specific Registry with the given alias.
// If the provided registry is nil a non-current Registry will be provided.
func NewAlias(alias string, registry Registry) *Alias {
	if registry == nil {
		registry = NewRegistry()
	}
	return &Alias{
		alias:    alias,
		Registry: registry,
	}
}

// Register the type for the specified example object.
// Generates the embedded Registry.Alias() call with first use.
// Actual registration passed along to package registry object.
func (a *Alias) Register(example interface{}) error {
	if !a.aliased {
		fmt.Println("nar", a.alias)
		if err := a.Alias(a.alias, example); err != nil {
			return fmt.Errorf("register alias: %w", err)
		}
		a.aliased = true
	}

	if err := a.Registry.Register(example); err != nil {
		return fmt.Errorf("register example: %w", err)
	}

	return nil
}
