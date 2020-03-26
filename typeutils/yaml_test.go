package typeutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

// These tests demonstrates and validates use of a Registry to marshal/unmarshal YAML.

type YamlTestSuite struct {
	suite.Suite
	film *filmYaml
}

var testRegistryYaml = NewRegistry()

func init() {
	if err := testRegistryYaml.Alias("test", filmYaml{}); err != nil {
		fmt.Printf("*** Error creating alias: %s\n", err)
	}
	if err := testRegistryYaml.Register(&alpha{}); err != nil {
		fmt.Printf("*** Error registering alpha: %s\n", err)
	}
	if err := testRegistryYaml.Register(&bravo{}); err != nil {
		fmt.Printf("*** Error registering bravo: %s\n", err)
	}
}

func (suite *YamlTestSuite) SetupTest() {
	copyFn = copyViaYaml
	suite.film = &filmYaml{Name: "Test YAML", Index: make(map[string]actor)}
	suite.film.Lead = &alpha{Name: "Goober", Percent: 13.23}
	suite.film.addActor("Goober", suite.film.Lead)
	suite.film.addActor("Snoofus", &bravo{Finished: false, Iterations: 17, extra: "stuff"})
	suite.film.addActor("Noodle", &alpha{Name: "Noodle", Percent: 19.57, extra: "stuff"})
	suite.film.addActor("Soup", &bravo{Finished: true, Iterations: 79})
}

func (suite *YamlTestSuite) TearDownSuite() {
	copyFn = nil
}

func TestYamlSuite(t *testing.T) {
	suite.Run(t, new(YamlTestSuite))
}

//////////////////////////////////////////////////////////////////////////

type filmYaml struct {
	yaml.Marshaler
	yaml.Unmarshaler

	Name  string
	Lead  actor
	Cast  []actor
	Index map[string]actor
}

type filmYamlConvert struct {
	Name  string
	Lead  interface{}
	Cast  []interface{}
	Index map[string]interface{}
}

func (film *filmYaml) addActor(name string, act actor) {
	film.Cast = append(film.Cast, act)
	film.Index[name] = act
}

func (film *filmYaml) MarshalYAML() (interface{}, error) {
	var err error

	convert := filmYamlConvert{
		Name: film.Name,
	}

	if convert.Lead, err = testRegistryYaml.ConvertItemToMap(film.Lead); err != nil {
		return nil, fmt.Errorf("converting lead to map: %w", err)
	}

	convert.Cast = make([]interface{}, len(film.Cast))
	for i, member := range film.Cast {
		if convert.Cast[i], err = testRegistryYaml.ConvertItemToMap(member); err != nil {
			return nil, fmt.Errorf("converting cast member to map: %w", err)
		}
	}

	convert.Index = make(map[string]interface{}, len(film.Index))
	for key, member := range film.Index {
		if convert.Index[key], err = testRegistryYaml.ConvertItemToMap(member); err != nil {
			return nil, fmt.Errorf("converting cast member to map: %w", err)
		}
	}

	return convert, nil
}

func (film *filmYaml) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("not a mapping node for filmYaml")
	}

	convert := filmYamlConvert{}
	if err := value.Decode(&convert); err != nil {
		return fmt.Errorf("decoding value to conversion struct: %w", err)
	}

	film.Name = convert.Name

	for i := 0; i < len(value.Content); i += 2 {
		var err error
		key := value.Content[i].Value
		val := value.Content[i+1]
		switch key {
		case "lead":
			if film.Lead, err = film.unmarshalActor(val); err != nil {
				return err
			}
		case "cast":
			film.Cast = make([]actor, len(val.Content))
			for i, node := range val.Content {
				if film.Cast[i], err = film.unmarshalActor(node); err != nil {
					return err
				}
			}
		case "index":
			film.Index = make(map[string]actor, len(val.Content)/2)
			for i := 0; i < len(val.Content); i += 2 {
				if film.Index[val.Content[i].Value], err = film.unmarshalActor(val.Content[i+1]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (film *filmYaml) unmarshalActor(value *yaml.Node) (actor, error) {
	if value.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("not a mapping node for actor")
	}

	var mapped map[string]interface{}
	if err := value.Decode(&mapped); err != nil {
		return nil, fmt.Errorf("decoding to map: %w", err)
	}

	if item, err := testRegistryYaml.CreateItemFromMap(mapped); err != nil {
		return nil, fmt.Errorf("creating item from map %v: %w", mapped, err)
	} else if act, ok := item.(actor); !ok {
		return nil, fmt.Errorf("item %v is not an actor", item)
	} else {
		return act, nil
	}
}

func copyViaYaml(dest, src interface{}) error {
	if bytes, err := yaml.Marshal(src); err != nil {
		return fmt.Errorf("marshaling from %v: %w", src, err)
	} else if err = yaml.Unmarshal(bytes, dest); err != nil {
		return fmt.Errorf("unmarshaling to %v: %w", dest, err)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////

// TestExample runs an example from some doc on yaml.v3 or something.
// It's here to taunt me when I can't figure out how it all works.
// Not directly applicable to this test suite.
func (suite *YamlTestSuite) TestExample() {
	type T struct {
		F int `yaml:"a,omitempty"`
		B int
	}
	t := T{F: 1, B: 2}
	bytes, err := yaml.Marshal(t)
	suite.Assert().NoError(err)
	var x T
	suite.Assert().NoError(yaml.Unmarshal(bytes, &x))
	suite.Assert().Equal(t, x)
}

func (suite *YamlTestSuite) TestCycle() {
	bytes, err := yaml.Marshal(suite.film)
	suite.Assert().NoError(err)

	//fmt.Printf(">>> marshaled:\n%s\n", string(bytes))

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
