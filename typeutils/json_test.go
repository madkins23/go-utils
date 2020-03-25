package typeutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type JsonTestSuite struct {
	suite.Suite
	film *filmJson
}

var testRegistryJson = NewRegistry()

func init() {
	if err := testRegistryJson.Alias("test", filmJson{}); err != nil {
		fmt.Printf("*** Error creating alias: %s\n", err)
	}
	if err := testRegistryJson.Register(&alpha{}); err != nil {
		fmt.Printf("*** Error registering alpha: %s\n", err)
	}
	if err := testRegistryJson.Register(&bravo{}); err != nil {
		fmt.Printf("*** Error registering bravo: %s\n", err)
	}
}

func (suite *JsonTestSuite) SetupTest() {
	suite.film = &filmJson{Name: "Test JSON", Index: make(map[string]actor)}
	suite.film.Lead = &alpha{Name: "Goober", Percent: 13.23}
	suite.film.addActor("Goober", suite.film.Lead)
	suite.film.addActor("Snoofus", &bravo{Finished: false, Iterations: 17, extra: "stuff"})
	suite.film.addActor("Noodle", &alpha{Name: "Noodle", Percent: 19.57, extra: "stuff"})
	suite.film.addActor("Soup", &bravo{Finished: true, Iterations: 79})
}

func TestJsonSuite(t *testing.T) {
	suite.Run(t, new(JsonTestSuite))
}

//////////////////////////////////////////////////////////////////////////

type filmJson struct {
	json.Marshaler
	json.Unmarshaler

	Name  string
	Lead  actor
	Cast  []actor
	Index map[string]actor
}

func (film *filmJson) addActor(name string, act actor) {
	film.Cast = append(film.Cast, act)
	film.Index[name] = act
}

func (film *filmJson) MarshalJSON() ([]byte, error) {
	var err error

	result := make(map[string]interface{})
	if result["lead"], err = testRegistryJson.ItemToMap(film.Lead); err != nil {
		return nil, fmt.Errorf("unable to convert lead to map: %w", err)
	}

	cast := make([]map[string]interface{}, len(film.Cast))
	for i, member := range film.Cast {
		if cast[i], err = testRegistryJson.ItemToMap(member); err != nil {
			return nil, fmt.Errorf("unable to convert cast member to map: %w", err)
		}
	}
	result["cast"] = cast

	index := make(map[string]map[string]interface{}, len(film.Index))
	for key, member := range film.Index {
		if index[key], err = testRegistryJson.ItemToMap(member); err != nil {
			return nil, fmt.Errorf("unable to convert cast member to map: %w", err)
		}
	}
	result["index"] = index

	return json.Marshal(result)
}

func (film *filmJson) UnmarshalJSON(input []byte) error {
	type filmJsonStruct struct {
	}

	var parse filmJsonStruct
	if err := json.Unmarshal(input, parse); err != nil {
		return fmt.Errorf("unable to unmarshal input JSON into struct: %w", err)
	}

	/*
		if value.Kind != yaml.MappingNode {
			return fmt.Errorf("not a mapping node for filmYaml")
		}

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
			default:
				return fmt.Errorf("unexpected film field: %s", key)
			}
		}
	*/

	return nil
}

func (film *filmJson) unmarshalActor(value *yaml.Node) (actor, error) {
	if value.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("not a mapping node for actor")
	}

	var err error
	var item interface{}
	for i := 0; i < len(value.Content); i += 2 {
		if value.Content[i].Value == TypeField {
			typeName := value.Content[i+1].Value
			if item, err = testRegistryJson.Make(typeName); err != nil {
				return nil, fmt.Errorf("unable to make registry instance %s: %w", typeName, err)
			}
		}
	}
	if item == nil {
		return nil, fmt.Errorf("no type field for actor node")
	}

	var act actor
	var ok bool
	if act, ok = item.(actor); !ok {
		return nil, fmt.Errorf("unable to convert %v to actor", reflect.TypeOf(item).Name())
	}

	if err = value.Decode(act); err != nil {
		return nil, fmt.Errorf("unable to decode actor: %w", err)
	}

	return act, nil
}

//////////////////////////////////////////////////////////////////////////

// TestExample runs an example from some doc on yaml.v3 or something.
// It's here to taunt me when I can't figure out how it all works.
func (suite *JsonTestSuite) TestExample() {
	type T struct {
		F int `json:"a,omitempty"`
		B int
	}
	t := T{F: 1, B: 2}
	bytes, err := json.Marshal(t)
	suite.Assert().NoError(err)
	var x T
	suite.Assert().NoError(json.Unmarshal(bytes, &x))
	suite.Assert().Equal(t, x)
}

func (suite *JsonTestSuite) TestMarshal() {
	/*
		bytes, err := yaml.Marshal(suite.film)
		suite.Assert().NoError(err)
		fmt.Println(string(bytes))
		var film filmJson
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
	*/
}
