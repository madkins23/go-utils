package typeutils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type YamlTestSuite struct {
	suite.Suite
	film *filmYaml
}

var testRegistry = NewRegistry()

func init() {
	if err := testRegistry.Alias("test", filmYaml{}); err != nil {
		fmt.Printf("*** Error creating alias: %s\n", err)
	}
	if err := testRegistry.Register(&alpha{}); err != nil {
		fmt.Printf("*** Error registering alpha: %s\n", err)
	}
	if err := testRegistry.Register(&bravo{}); err != nil {
		fmt.Printf("*** Error registering bravo: %s\n", err)
	}
}

func (suite *YamlTestSuite) SetupTest() {
	suite.film = newFilmYaml()
	suite.film.Lead = &alpha{Name: "Goober", Percent: 13.23}
	suite.film.addActor("Goober", suite.film.Lead)
	suite.film.addActor("Snoofus", &bravo{Finished: false, Iterations: 17, extra: "stuff"})
	suite.film.addActor("Noodle", &alpha{Name: "Noodle", Percent: 19.57, extra: "stuff"})
	suite.film.addActor("Soup", &bravo{Finished: true, Iterations: 79})
}

func TestYamlSuite(t *testing.T) {
	suite.Run(t, new(YamlTestSuite))
}

//////////////////////////////////////////////////////////////////////////

type filmYaml struct {
	yaml.Marshaler
	yaml.Unmarshaler

	Lead  actor
	Cast  []actor
	Index map[string]actor
}

func (dir *filmYaml) addActor(name string, act actor) {
	dir.Cast = append(dir.Cast, act)
	dir.Index[name] = act
}

func newFilmYaml() *filmYaml {
	return &filmYaml{Index: make(map[string]actor)}
}

func (dir *filmYaml) MarshalYAML() (interface{}, error) {
	var err error

	result := make(map[string]interface{})
	if result["lead"], err = testRegistry.ItemToMap(dir.Lead); err != nil {
		return nil, fmt.Errorf("unable to convert lead to map: %w", err)
	}

	cast := make([]map[string]interface{}, len(dir.Cast))
	for i, member := range dir.Cast {
		if cast[i], err = testRegistry.ItemToMap(member); err != nil {
			return nil, fmt.Errorf("unable to convert cast member to map: %w", err)
		}
	}
	result["cast"] = cast

	index := make(map[string]map[string]interface{}, len(dir.Index))
	for key, member := range dir.Index {
		if index[key], err = testRegistry.ItemToMap(member); err != nil {
			return nil, fmt.Errorf("unable to convert cast member to map: %w", err)
		}
	}
	result["index"] = index

	return result, nil
}

func unmarshalActor(value *yaml.Node) (actor, error) {
	if value.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("not a mapping node for actor")
	}

	var err error
	var mapped map[string]interface{}
	if err = value.Decode(&mapped); err != nil {
		return nil, fmt.Errorf("unable to decode actor map: %w", err)
	}

	var mud interface{}
	if mud, err = testRegistry.MapToItem(mapped); err != nil {
		return nil, fmt.Errorf("unable to map to actor: %w", err)
	}

	var act actor
	var ok bool
	if act, ok = mud.(actor); !ok {
		return nil, fmt.Errorf("unable to convert %v to actor", reflect.TypeOf(mud).Name())
	}

	if err = value.Decode(act); err != nil {
		return nil, fmt.Errorf("unable to decode actor: %w", err)
	}

	return act, nil
}

func (dir *filmYaml) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("not a mapping node for filmYaml")
	}

	for i := 0; i < len(value.Content); i += 2 {
		var err error
		key := value.Content[i].Value
		val := value.Content[i+1]
		switch key {
		case "lead":
			if dir.Lead, err = unmarshalActor(val); err != nil {
				return err
			}
		case "cast":
			dir.Cast = make([]actor, len(val.Content))
			for i, node := range val.Content {
				if dir.Cast[i], err = unmarshalActor(node); err != nil {
					return err
				}
			}
		case "index":
			dir.Index = make(map[string]actor, len(val.Content)/2)
			for i := 0; i < len(val.Content); i += 2 {
				if dir.Index[val.Content[i].Value], err = unmarshalActor(val.Content[i+1]); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unexpected film field: %s", key)
		}
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////

// TestExample runs an example from some doc on yaml.v3 or something.
// It's here to taunt me when I can't figure out how it all works.
func (suite *YamlTestSuite) TestExample() {
	type T struct {
		F int `yaml:"a,omitempty"`
		B int
	}
	var t T
	suite.Assert().NoError(yaml.Unmarshal([]byte("a: 1\nb: 2"), &t))
	suite.Assert().Equal(T{F: 1, B: 2}, t)
}

func (suite *YamlTestSuite) TestMarshal() {
	bytes, err := yaml.Marshal(suite.film)
	suite.Assert().NoError(err)
	fmt.Println(string(bytes))
	var film filmYaml
	suite.Assert().NoError(yaml.Unmarshal(bytes, &film))
	suite.Assert().NotEqual(suite.film, &film) // fails due to unexported field 'extra'
	for _, act := range suite.film.Cast {
		// Remove unexported field.
		if alf, ok := act.(*alpha); ok {
			alf.extra = ""
		} else if bra, ok := act.(*bravo); ok {
			bra.extra = ""
		}
	}
	suite.Assert().Equal(suite.film, &film) // succeeds now that unexported fields are gone.
}
