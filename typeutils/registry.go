package typeutils

import (
	"fmt"
	"reflect"
	"strings"
)

const TypeField = "$type$"

type FromMapFn func(from map[string]interface{}, to interface{}) error
type ToMapFn func(from interface{}, to map[string]interface{}) error

// Registry is the type registry interface.
// A type registry tracks specific types by name, a facility not native to Go.
// A type name in the registry is made up of package path and local type name.
// Aliases may be specified to shorten the path to manageable lengths.
type Registry interface {
	Alias(alias string, example interface{}) error
	Register(example interface{}) error
	Make(name string) (interface{}, error)
	NameFor(item interface{}) (string, error)
	GenNames(item interface{}, aliased bool) (string, []string, error)
	ConvertItemToMap(item interface{}) (map[string]interface{}, error)
	CreateItemFromMap(in map[string]interface{}) (interface{}, error)
}

// RegistryItem contains methods for pushing fields to or pulling fields from a map.
// A Registry will work with any kind of struct, but won't copy field data without this interface.
// This is used by ConvertItemToMap and CreateItemFromMap (as called from marshal/unmarshal code).
// Note that both methods must be provided for either to work.
type RegistryItem interface {
	PushToMap(toMap map[string]interface{}) error
	PullFromMap(fromMap map[string]interface{}) error
}

// NewRegistry creates a new Registry object of the default internal type.
// Registries created via this function are not safe for concurrent access,
// manage this access or use NewRegistrar() to create a concurrent safe version.
// Developers might be able to write more efficient concurrency code using Registry.
func NewRegistry() Registry {
	return &registry{
		byName:  make(map[string]*registration),
		byType:  make(map[reflect.Type]*registration),
		aliases: make(map[string]string),
	}
}

//////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////

// Default Registry implementation.
type registry struct {
	// byName supports lookup of registrations by 'name'.
	// Full names and aliases are both entered herein.
	byName map[string]*registration

	// byType supports lookup of registrations by type.
	byType map[reflect.Type]*registration

	// alias maps shortened 'alias' strings to path prefix to shorten names.
	aliases map[string]string
}

// Registration structure groups data from indexes.
type registration struct {
	// typeName includes package path and type name.
	defaultName string

	// typeNames is the set of all possible (i.e. aliased) type names.
	// The best one will always be in typeName.
	allNames []string

	// typeObj is the reflect.Type object for the example object.
	typeObj reflect.Type
}

//////////////////////////////////////////////////////////////////////////

// Alias creates an alias to be used to shorten names.
// Alias must exist prior to registering applicable types.
// Redefining a pre-existing alias is an error.
func (reg *registry) Alias(alias string, example interface{}) error {
	if _, ok := reg.aliases[alias]; ok {
		return fmt.Errorf("can't redefine alias %s", alias)
	}

	exampleType := reflect.TypeOf(example)
	if exampleType == nil {
		return fmt.Errorf("no type for alias %s (%v)", alias, example)
	}

	if exampleType.Kind() == reflect.Ptr {
		exampleType = exampleType.Elem()
		if exampleType == nil {
			return fmt.Errorf("no elem type for alias %s (%v)", alias, example)
		}
	}

	pkgPath := exampleType.PkgPath()
	if pkgPath == "" {
		return fmt.Errorf("no package path for alias %s (%v)", alias, example)
	}

	reg.aliases[alias] = pkgPath
	return nil
}

// Register a type by providing an example object.
func (reg *registry) Register(example interface{}) error {
	// Get reflected type for example object.
	exType := reflect.TypeOf(example)
	if exType != nil && exType.Kind() == reflect.Ptr {
		exType = exType.Elem()
	}
	if exType == nil {
		return fmt.Errorf("no reflected type for %v", example)
	}

	// Check for previous record.
	if _, ok := reg.byType[exType]; ok {
		return fmt.Errorf("previous registration for type %v", exType)
	}

	// Get type name without any pointer asterisks.
	typeName := exType.String()
	if strings.HasPrefix(typeName, "*") {
		typeName = strings.TrimLeft(typeName, "*")
	}

	// Create registration record for this type.
	item := &registration{
		defaultName: typeName,
		allNames:    make([]string, 1, len(reg.aliases)+1),
		typeObj:     exType,
	}

	// Initialize default name to full name with package and type.
	name, aliases, err := reg.GenNames(example, true)
	if err != nil {
		return fmt.Errorf("getting type name of example: %w", err)
	}

	item.defaultName = name
	item.allNames[0] = name
	for _, alias := range aliases {
		item.allNames = append(item.allNames, alias)
	}

	// Add name lookups for all default and aliased names.
	reg.byName[name] = item
	for _, name := range item.allNames {
		reg.byName[name] = item
	}

	// Add type lookup.
	reg.byType[exType] = item

	return nil
}

// GenNames creates the possible names for the type represented by the example object.
// Returns the 'canonical' name, an optional array of aliased names per current aliases, and any error.
// If the aliased argument is true a possibly empty array will be returned for the second argument otherwise nil.
func (reg *registry) GenNames(example interface{}, aliased bool) (string, []string, error) {
	// Initialize default name to full name with package and type.
	name, err := genNameFromInterface(example)
	if err != nil {
		return "", nil, fmt.Errorf("generating basic name: %w", err)
	}

	var aliases []string
	if aliased {
		aliases = make([]string, 0, len(reg.aliases))

		// Look for any possible aliases for the type and add them to the list of all names.
		for alias, prefixPath := range reg.aliases {
			if strings.HasPrefix(name, prefixPath) {
				aliases = append(aliases, "["+alias+"]"+name[len(prefixPath)+1:])
			}
		}

		// Choose default name again from shortest, therefore most likely an aliased name if there are any.
		nameLen := len(name)
		for _, alias := range aliases {
			// Using <= favors later aliases of same size.
			if len(alias) <= nameLen {
				name = alias
			}
		}
	}

	return name, aliases, nil
}

// NameFor returns a name for the specified object.
func (reg *registry) NameFor(item interface{}) (string, error) {
	itemType := reflect.TypeOf(item)
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	registration, ok := reg.byType[itemType]
	if !ok {
		return "", fmt.Errorf("no registration for type %s", itemType)
	}

	return registration.defaultName, nil
}

// Make creates a new instance of the example object with the specified name.
// The new instance will be created with fields filled with zero values.
func (reg *registry) Make(name string) (interface{}, error) {
	item, found := reg.byName[name]
	if !found {
		return nil, fmt.Errorf("no example for '%s'", name)
	}

	return reflect.New(item.typeObj).Interface(), nil
}

// ConvertItemToMap converts a registry typed item into a map for further processing.
// If the item is not of a Registry type an error is returned.
func (reg *registry) ConvertItemToMap(item interface{}) (map[string]interface{}, error) {
	value := reflect.ValueOf(item)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("item %s is not a struct", item)
	}

	itemType := value.Type()
	registration, ok := reg.byType[itemType]
	if !ok {
		return nil, fmt.Errorf("no registration for type %s", itemType)
	}

	result := make(map[string]interface{})

	// Add the special marker for the type of the object.
	// This should work with both JSON and YAML.
	result[TypeField] = registration.defaultName

	if mapper, ok := item.(RegistryItem); ok {
		if err := mapper.PushToMap(result); err != nil {
			return nil, fmt.Errorf("pushing item fields to map: %w", err)
		}
	}

	return result, nil
}

// CreateItemFromMap attempts to return a new item of the type specified in the map.
// An error is returned if this is impossible.
func (reg *registry) CreateItemFromMap(in map[string]interface{}) (interface{}, error) {
	typeField, found := in[TypeField]
	if !found {
		_ = fmt.Errorf("no object type in map")
	}
	typeName, ok := typeField.(string)
	if !ok {
		_ = fmt.Errorf("converting type field %v to string", typeField)
	}

	item, err := reg.Make(typeName)
	if err != nil {
		return nil, fmt.Errorf("making item of type %s: %w", typeField, err)
	}

	if mapper, ok := item.(RegistryItem); ok {
		if err := mapper.PullFromMap(in); err != nil {
			return nil, fmt.Errorf("pulling item fields from map: %w", err)
		}
	}

	return item, nil
}

//////////////////////////////////////////////////////////////////////////

func genNameFromInterface(example interface{}) (string, error) {
	itemType := reflect.TypeOf(example)
	if itemType == nil {
		return "", fmt.Errorf("no type for item %v", example)
	}

	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	path := itemType.PkgPath()
	if path == "" {
		return "", fmt.Errorf("no path for type %s", itemType.Name())
	}

	last := strings.LastIndex(path, "/")
	if last < 0 {
		return "", fmt.Errorf("no slash in %s", path)
	}

	final := path[last:]
	name := itemType.Name()

	if strings.HasPrefix(name, final+".") {
		name = name[len(final)+1:]
	}

	return path + "/" + name, nil
}
