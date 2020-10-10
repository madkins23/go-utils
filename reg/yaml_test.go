package reg

import (
	"fmt"
	"io/ioutil"
	"strings"
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
	if err := testRegistryYaml.Register(&filmYaml{}); err != nil {
		fmt.Printf("*** Error registering alpha: %s\n", err)
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

func (suite *YamlTestSuite) TestYamlBaseIsRegistry() {
	var reg Registry = &YamlBase{}
	_, ok := reg.(Registry)
	suite.Assert().True(ok)
}

func (suite *YamlTestSuite) TestGetTypeNameAndReset() {
	reader := strings.NewReader(simpleYaml)
	suite.Assert().NotNil(reader)

	// Get type name.
	name, err := getYamlTypeNameAndReset(reader)
	suite.Assert().NoError(err)
	suite.Assert().Equal("[test]filmYaml", name)

	// Check for reset.
	bytes, err := ioutil.ReadAll(reader)
	suite.Assert().NoError(err)
	suite.Assert().Equal(simpleYaml, string(bytes))
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

func (film *filmYaml) addActor(name string, act actor) {
	film.Cast = append(film.Cast, act)
	film.Index[name] = act
}

func (film *filmYaml) PullFromMap(fromMap map[string]interface{}) error {
	if name, found := fromMap["name"].(string); found {
		film.Name = name
	}

	var ok bool
	if fromMap["name"] != nil {
		if film.Name, ok = fromMap["name"].(string); !ok {
			return fmt.Errorf("film name is not a string")
		}
	}

	var err error
	if fromMap["lead"] != nil {
		if film.Lead, err = film.pullActorFromMap(fromMap["lead"]); err != nil {
			return fmt.Errorf("pull lead actor from map: %w", err)
		}
	}

	if castElement, found := fromMap["cast"]; found && castElement != nil {
		if cast, ok := castElement.([]interface{}); ok {
			film.Cast = make([]actor, 0, len(cast))
			for _, actMap := range cast {
				if act, err := film.pullActorFromMap(actMap); err != nil {
					return fmt.Errorf("pulling actor from map: %w", err)
				} else {
					film.Cast = append(film.Cast, act)
				}
			}
		}
	}

	if indexElement, found := fromMap["index"]; found && indexElement != nil {
		if index, ok := indexElement.(map[string]interface{}); ok {
			film.Index = make(map[string]actor)
			for key, actMap := range index {
				if act, err := film.pullActorFromMap(actMap); err != nil {
					return fmt.Errorf("pulling actor from map: %w", err)
				} else {
					film.Index[key] = act
				}
			}
		}
	}

	return nil
}

func (film *filmYaml) pullActorFromMap(from interface{}) (actor, error) {
	if fromMap, ok := from.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("from is not a map")
	} else if actItem, err := testRegistryYaml.CreateItemFromMap(fromMap); err != nil {
		return nil, fmt.Errorf("create actor from map: %w", err)
	} else if act, ok := actItem.(actor); !ok {
		return nil, fmt.Errorf("map is not an actor")
	} else {
		return act, nil
	}
}

func (film *filmYaml) PushToMap(toMap map[string]interface{}) error {
	var err error
	if toMap[TypeField], err = testRegistryYaml.NameFor(film); err != nil {
		return fmt.Errorf("get type name: %w", err)
	}

	toMap["name"] = film.Name

	if toMap["lead"], err = testRegistryYaml.ConvertItemToMap(film.Lead); err != nil {
		return fmt.Errorf("converting lead to map: %w", err)
	}

	cast := make([]interface{}, len(film.Cast))
	for i, member := range film.Cast {
		if cast[i], err = testRegistryYaml.ConvertItemToMap(member); err != nil {
			return fmt.Errorf("converting cast member to map: %w", err)
		}
	}
	toMap["cast"] = cast

	index := make(map[string]interface{}, len(film.Index))
	for key, member := range film.Index {
		if index[key], err = testRegistryYaml.ConvertItemToMap(member); err != nil {
			return fmt.Errorf("converting cast member to map: %w", err)
		}
	}
	toMap["index"] = index

	return nil
}

func (film *filmYaml) MarshalYAML() (interface{}, error) {
	toMap := make(map[string]interface{})
	if err := film.PushToMap(toMap); err != nil {
		return nil, fmt.Errorf("pushing film to map: %w", err)
	} else {
		return toMap, nil
	}
}

func (film *filmYaml) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("not a mapping node for film")
	}

	// Simpler to code than pulling everything bit by bit from the value object.
	// The latter might be faster, however.
	temp := make(map[string]interface{})
	if err := value.Decode(temp); err != nil {
		return fmt.Errorf("decoding film to temp: %w", err)
	}

	if err := film.PullFromMap(temp); err != nil {
		return fmt.Errorf("pulling film from map: %w", err)
	} else {
		return nil
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

func (suite *YamlTestSuite) TestMarshalCycle() {
	bytes, err := yaml.Marshal(suite.film)
	suite.Assert().NoError(err)

	fmt.Print("--- TestMarshalCycle ---------------\n", string(bytes), "------------------------------------\n")

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

func (suite *YamlTestSuite) TestLoadFromString() {
	base := NewYamlBase(testRegistryYaml)
	loaded, err := base.LoadFromString(simpleYaml)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(loaded)
	film, ok := loaded.(*filmYaml)
	suite.Assert().True(ok)
	suite.Assert().NotNil(film)
	suite.Assert().Equal("Blockbuster Movie", film.Name)
	suite.checkAlpha(film.Lead)
	suite.Assert().NotNil(film.Cast)
	suite.Assert().Len(film.Cast, 2)
	suite.checkAlpha(film.Cast[0])
	suite.checkBravo(film.Cast[1])
	suite.Assert().NotNil(film.Index)
	suite.Assert().Len(film.Index, 2)
	suite.checkAlpha(film.Index["Lucky, Lance"])
	suite.checkBravo(film.Index["Queue, Susie"])
}

func (suite *YamlTestSuite) TestMarshalFileCycle() {
	base := NewYamlBase(testRegistryYaml)
	file, err := ioutil.TempFile("", "*.test.yaml")
	suite.Assert().NoError(err)
	suite.Assert().NotNil(file)
	fileName := file.Name()
	// Go ahead and close it, just needed the file name.
	suite.Assert().NoError(file.Close())

	film := suite.makeTestFilm()
	suite.Assert().NoError(base.SaveToFile(film, fileName))

	reloaded, err := base.LoadFromFile(fileName)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reloaded)
	suite.Assert().Equal(film, reloaded)
}
func (suite *YamlTestSuite) TestMarshalStringCycle() {
	base := NewYamlBase(testRegistryYaml)
	film := suite.makeTestFilm()
	str, err := base.SaveToString(film)
	suite.Assert().NoError(err)
	suite.NotZero(str)

	fmt.Print("--- TestMarshalStringCycle ---------\n", str, "------------------------------------\n")

	reloaded, err := base.LoadFromString(str)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reloaded)
	suite.Assert().Equal(film, reloaded)
}

//////////////////////////////////////////////////////////////////////////

const simpleYaml = `
<type>: '[test]filmYaml'
name:   'Blockbuster Movie'
lead: {
  <type>: '[test]alpha',
  name: 'Lance Lucky',
  percentDone: 23.79,
  extra: 'Yaaaa!'
}
cast:
- {
    <type>: '[test]alpha',
    name: 'Lance Lucky',
    percentDone: 23.79,
    extra: false
  }
- {
    <type>: '[test]bravo',
    finished: true,
    iterations: 13,
    extra: 'gibbering ghostwhistle'
  }
index: {
  'Lucky, Lance': {
    <type>: '[test]alpha',
    name: 'Lance Lucky',
    percentDone: 23.79,
    extra: 'marshmallow stars'
  },
  'Queue, Susie': {
    <type>: '[test]bravo',
    finished: true,
    iterations: 13,
    extra: 19.57
  }
}
`

func (suite *YamlTestSuite) checkAlpha(act actor) {
	suite.Assert().NotNil(act)
	alf, ok := act.(*alpha)
	suite.Assert().True(ok)
	suite.Assert().NotNil(alf)
	suite.Assert().Equal("Lance Lucky", alf.Name)
	suite.Assert().Equal(float32(23.79), alf.Percent)
	suite.Assert().Empty(alf.extra)
}

func (suite *YamlTestSuite) checkBravo(act actor) {
	suite.Assert().NotNil(act)
	bra, ok := act.(*bravo)
	suite.Assert().True(ok)
	suite.Assert().NotNil(bra)
	suite.Assert().True(bra.Finished)
	suite.Assert().Equal(13, bra.Iterations)
	suite.Assert().Empty(bra.extra)
}

func (suite *YamlTestSuite) makeTestFilm() *filmYaml {
	actor1 := &alpha{
		Name:    "Goober Snoofus",
		Percent: 13.23,
	}
	actor2 := &bravo{
		Finished:   true,
		Iterations: 1957,
	}
	return &filmYaml{
		Name: "",
		Lead: actor1,
		Cast: []actor{
			actor1,
			actor2,
		},
		Index: map[string]actor{
			"Snoofus, Goober": actor1,
			"Snarly, Booger":  actor2,
		},
	}
}
