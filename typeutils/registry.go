package typeutils

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

const TypeField = "$type$"

// Registry is the type registry interface.
// A type registry tracks specific types by name, a facility not native to Go.
// A type name in the registry is made up of package path and local type name.
// Aliases may be specified to shorten the path to manageable lengths.
type Registry interface {
	Alias(alias string, example interface{}) error
	Register(example interface{}) error
	Make(name string) (interface{}, error)
	NameFor(item interface{}) (string, error)
	ItemToMap(item interface{}) (map[string]interface{}, error)
	MapToItem(map[string]interface{}) (interface{}, error)
}

// NewRegistry creates a new Registry object of the default internal type.
// Registries created via this function are not safe for concurrent access,
// manage this access or use NewRegistrar() to create a concurrent safe version.
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
// Use an empty string to remove a previous alias.
func (reg *registry) Alias(alias string, example interface{}) error {
	exampleType := reflect.TypeOf(example)
	if exampleType == nil {
		return fmt.Errorf("no type for example %v", example)
	}

	pkgPath := exampleType.PkgPath()
	if pkgPath != "" {
		reg.aliases[alias] = pkgPath
	} else if _, ok := reg.aliases[alias]; ok {
		delete(reg.aliases, alias)
	}

	// TODO: change all registration records to reflect change in alias?
	// Could just assume all aliases are in place prior to registration.

	return nil
}

// Register a type by providing an example object.
func (reg *registry) Register(example interface{}) error {
	// Get reflected type for example object.
	exampleType := reflect.TypeOf(example)
	if exampleType == nil {
		return fmt.Errorf("no reflected type for %v", example)
	}
	if exampleType.Kind() == reflect.Ptr {
		exampleType = exampleType.Elem()
	}

	// Check for previous record.
	if _, ok := reg.byType[exampleType]; ok {
		return fmt.Errorf("previous registration for type %v", exampleType)
	}

	// Get type name without any pointer asterisks.
	typeName := exampleType.String()
	if strings.HasPrefix(typeName, "*") {
		typeName = strings.TrimLeft(typeName, "*")
	}

	// Create registration record for this type.
	item := &registration{
		defaultName: typeName,
		allNames:    make([]string, 1, len(reg.aliases)+1),
		typeObj:     exampleType,
	}

	// Initialize default name to full name with package and type.
	var err error
	item.defaultName, err = genNameFromInterface(example)
	if err != nil {
		return err
	}
	item.allNames[0] = item.defaultName

	// Look for any possible aliases for the type and add them to the list of all names.
	for alias, prefixPath := range reg.aliases {
		if strings.HasPrefix(item.defaultName, prefixPath) {
			aliasedName := "[" + alias + "]" + item.defaultName[len(prefixPath)+1:]
			item.allNames = append(item.allNames, aliasedName)
		}
	}

	// Choose default name again from shortest, therefore most likely an aliased name if there are any.
	for _, name := range item.allNames[1:] {
		// Using >= favors later aliases of same size.
		if len(name) <= len(item.defaultName) {
			item.defaultName = name
		}
	}

	// Add name lookups for all default and aliased names.
	for _, name := range item.allNames {
		reg.byName[name] = item
	}

	// Add type lookup.
	reg.byType[exampleType] = item

	return nil
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

// NameFor returns a name for the specified object.
func (reg *registry) NameFor(item interface{}) (string, error) {
	itemType := reflect.TypeOf(item)
	registration, ok := reg.byType[itemType]
	if !ok {
		return "", fmt.Errorf("no registration for type %s", itemType)
	}

	return registration.defaultName, nil
}

// marshalYAML converts a registry typed item into a map for further processing.
// If the item is not of a Registry type an error is returned.
func (reg *registry) ItemToMap(item interface{}) (map[string]interface{}, error) {
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

	if err := copyFieldsToMap(value, result); err != nil {
		return nil, fmt.Errorf(": %w", err)
	}

	return result, nil
}

func copyFieldsToMap(value reflect.Value, result map[string]interface{}) error {
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)

		if field.Anonymous {
			if subValue := value.Field(i); subValue.Type().Kind() == reflect.Struct {
				// Process anonymous sub-structs as if they were part of this item.
				if err := copyFieldsToMap(subValue, result); err != nil {
					return err
				}
			}
		} else if unicode.IsUpper(rune(field.Name[0])) {
			fieldName := field.Name
			tagValue := field.Tag.Get("yaml")
			if tagValue == "" {
				fieldName = strings.ToLower(fieldName)
			} else if tagFields := strings.Split(tagValue, ","); len(tagFields) > 0 {
				if tagFields[0] != "" {
					fieldName = tagFields[0]
				}

				// TODO: other tag fields?
			}

			fld := value.Field(i)
			if fld.Interface() != nil {
				result[fieldName] = fld.Interface()
			}
		}
	}

	return nil
}

// MapToItem attempts to return a new item of the type specified in the map.
// An error is returned if this is impossible.
func (reg *registry) MapToItem(in map[string]interface{}) (interface{}, error) {
	typeField, found := in[TypeField]
	if !found {
		_ = fmt.Errorf("unable to find type for object")
	}
	typeName, ok := typeField.(string)
	if !ok {
		_ = fmt.Errorf("unable to convert type field %v to string", typeField)
	}

	item, err := reg.Make(typeName)
	if err != nil {
		return nil, fmt.Errorf("unable to make item of type %s", typeField)
	}

	// TODO: should initialization from map occur here?
	// Apparently not, since this would require knowing about JSON/YAML here.
	// Could, however, use a function pointer,

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
