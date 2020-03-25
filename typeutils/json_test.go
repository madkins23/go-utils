package typeutils

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
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
	json.Marshaler   `json:",omitempty"`
	json.Unmarshaler `json:",omitempty"`

	Name  string
	Lead  actor
	Cast  []actor
	Index map[string]actor
}

type filmJsonConvert struct {
	Name  string
	Lead  interface{}
	Cast  []interface{}
	Index map[string]interface{}
}

func (film *filmJson) addActor(name string, act actor) {
	film.Cast = append(film.Cast, act)
	film.Index[name] = act
}

func (film *filmJson) MarshalJSON() ([]byte, error) {
	var err error

	convert := filmJsonConvert{
		Name: film.Name,
	}

	if convert.Lead, err = testRegistryJson.ItemToMap(film.Lead, toMapJson); err != nil {
		return nil, fmt.Errorf("unable to convert lead to map: %w", err)
	}

	convert.Cast = make([]interface{}, len(film.Cast))
	for i, member := range film.Cast {
		if convert.Cast[i], err = testRegistryJson.ItemToMap(member, toMapJson); err != nil {
			return nil, fmt.Errorf("unable to convert cast member to map: %w", err)
		}
	}

	convert.Index = make(map[string]interface{}, len(film.Index))
	for key, member := range film.Index {
		if convert.Index[key], err = testRegistryJson.ItemToMap(member, toMapJson); err != nil {
			return nil, fmt.Errorf("unable to convert cast member to map: %w", err)
		}
	}

	return json.Marshal(convert)
}

func (film *filmJson) UnmarshalJSON(input []byte) error {
	var err error

	convert := filmJsonConvert{}
	if err = json.Unmarshal(input, &convert); err != nil {
		return fmt.Errorf("unable to unmarshal input JSON into struct: %w", err)
	}

	film.Name = convert.Name

	if film.Lead, err = film.unmarshalActor(convert.Lead); err != nil {
		return fmt.Errorf("unable to unmarshal lead actor: %w", err)
	}

	film.Cast = make([]actor, len(convert.Cast))
	for i, member := range convert.Cast {
		if film.Cast[i], err = film.unmarshalActor(member); err != nil {
			return fmt.Errorf("unable to unmarshal cast member: %w", err)
		}
	}

	film.Index = make(map[string]actor, len(convert.Index))
	for name, member := range convert.Index {
		if film.Index[name], err = film.unmarshalActor(member); err != nil {
			return fmt.Errorf("unable to unmarshal index member: %w", err)
		}
	}

	return nil
}

func (film *filmJson) unmarshalActor(input interface{}) (actor, error) {
	actMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("actor input should be map")
	}

	if item, err := testRegistryJson.MapToItem(actMap, fromMapJson); err != nil {
		return nil, fmt.Errorf("unable to map to item: %w", err)
	} else if act, ok := item.(actor); !ok {
		return nil, fmt.Errorf("item is not an actor")
	} else {
		return act, nil
	}
}

func fromMapJson(from map[string]interface{}, to interface{}) error {
	if bytes, err := json.Marshal(from); err != nil {
		return fmt.Errorf("unable to marshal from %v: %w", from, err)
	} else if err = json.Unmarshal(bytes, &to); err != nil {
		return fmt.Errorf("unable to unmarshal to %v: %w", to, err)
	}

	return nil
}

func toMapJson(from interface{}, to map[string]interface{}) error {
	if bytes, err := json.Marshal(from); err != nil {
		return fmt.Errorf("unable to marshal from %v: %w", from, err)
	} else if err = json.Unmarshal(bytes, &to); err != nil {
		return fmt.Errorf("unable to unmarshal to %v: %w", to, err)
	}

	return nil
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
	bytes, err := json.Marshal(suite.film)
	suite.Assert().NoError(err)
	fmt.Println(string(bytes))
	var film filmJson
	suite.Assert().NoError(json.Unmarshal(bytes, &film))
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
