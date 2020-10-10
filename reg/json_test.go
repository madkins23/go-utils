package reg

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// These tests demonstrates and validates use of a Registry to marshal/unmarshal JSON.

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
	copyMapFromItemFn = copyItemToMap
	copyItemFromMapFn = copyMapToItem
	suite.film = &filmJson{Name: "Test JSON", Index: make(map[string]actor)}
	suite.film.Lead = &alpha{Name: "Goober", Percent: 13.23}
	suite.film.addActor("Goober", suite.film.Lead)
	suite.film.addActor("Snoofus", &bravo{Finished: false, Iterations: 17, extra: "stuff"})
	suite.film.addActor("Noodle", &alpha{Name: "Noodle", Percent: 19.57, extra: "stuff"})
	suite.film.addActor("Soup", &bravo{Finished: true, Iterations: 79})
}

func (suite *JsonTestSuite) TearDownSuite() {
	copyMapFromItemFn = nil
	copyItemFromMapFn = nil
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

	if convert.Lead, err = testRegistryJson.ConvertItemToMap(film.Lead); err != nil {
		return nil, fmt.Errorf("converting lead to map: %w", err)
	}

	convert.Cast = make([]interface{}, len(film.Cast))
	for i, member := range film.Cast {
		if convert.Cast[i], err = testRegistryJson.ConvertItemToMap(member); err != nil {
			return nil, fmt.Errorf("converting cast member to map: %w", err)
		}
	}

	convert.Index = make(map[string]interface{}, len(film.Index))
	for key, member := range film.Index {
		if convert.Index[key], err = testRegistryJson.ConvertItemToMap(member); err != nil {
			return nil, fmt.Errorf("converting cast member to map: %w", err)
		}
	}

	return json.Marshal(convert)
}

func (film *filmJson) UnmarshalJSON(input []byte) error {
	var err error

	convert := filmJsonConvert{}
	if err = json.Unmarshal(input, &convert); err != nil {
		return fmt.Errorf("unmarshaling input JSON into struct: %w", err)
	}

	film.Name = convert.Name

	if film.Lead, err = film.unmarshalActor(convert.Lead); err != nil {
		return fmt.Errorf("unmarshaling lead actor: %w", err)
	}

	film.Cast = make([]actor, len(convert.Cast))
	for i, member := range convert.Cast {
		if film.Cast[i], err = film.unmarshalActor(member); err != nil {
			return fmt.Errorf("unmarshaling cast member: %w", err)
		}
	}

	film.Index = make(map[string]actor, len(convert.Index))
	for name, member := range convert.Index {
		if film.Index[name], err = film.unmarshalActor(member); err != nil {
			return fmt.Errorf("unmarshaling index member: %w", err)
		}
	}

	return nil
}

func (film *filmJson) unmarshalActor(input interface{}) (actor, error) {
	actMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("actor input should be map")
	}

	if item, err := testRegistryJson.CreateItemFromMap(actMap); err != nil {
		return nil, fmt.Errorf("creating item from map: %w", err)
	} else if act, ok := item.(actor); !ok {
		return nil, fmt.Errorf("item is not an actor")
	} else {
		return act, nil
	}
}

func copyItemToMap(toMap map[string]interface{}, fromItem interface{}) error {
	if bytes, err := json.Marshal(fromItem); err != nil {
		return fmt.Errorf("marshaling from %v: %w", fromItem, err)
	} else if err = json.Unmarshal(bytes, &toMap); err != nil {
		return fmt.Errorf("unmarshaling to %v: %w", toMap, err)
	}

	return nil
}

func copyMapToItem(toItem interface{}, fromMap map[string]interface{}) error {
	if bytes, err := json.Marshal(fromMap); err != nil {
		return fmt.Errorf("marshaling from %v: %w", fromMap, err)
	} else if err = json.Unmarshal(bytes, toItem); err != nil {
		return fmt.Errorf("unmarshaling to %v: %w", toItem, err)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////

// TestExample duplicates the YAML test.
// Not directly applicable to this test suite.
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

func (suite *JsonTestSuite) TestCycle() {
	bytes, err := json.Marshal(suite.film)
	suite.Assert().NoError(err)

	//fmt.Printf(">>> marshaled:\n%s\n", string(bytes))

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
