package typeutils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

//////////////////////////////////////////////////////////////////////////

type RegistryTestSuite struct {
	suite.Suite
	registry Registry
	reg      *registry
}

func (suite *RegistryTestSuite) SetupTest() {
	// See json_test.go for these definitions.
	copyMapFromItemFn = copyItemToMap
	copyItemFromMapFn = copyMapToItem

	suite.registry = NewRegistry()
	var ok bool
	suite.reg, ok = suite.registry.(*registry)
	suite.Assert().True(ok)
}

func (suite *RegistryTestSuite) TearDownSuite() {
	copyMapFromItemFn = nil
	copyItemFromMapFn = nil
}

func TestRegistrySuite(t *testing.T) {
	suite.Run(t, new(RegistryTestSuite))
}

//////////////////////////////////////////////////////////////////////////

func (suite *RegistryTestSuite) TestNewRegistry() {
	suite.Assert().NotNil(suite.registry)
	suite.Assert().NotNil(suite.reg.byName)
	suite.Assert().NotNil(suite.reg.byType)
}

func (suite *RegistryTestSuite) TestAlias() {
	example := &alpha{}
	err := suite.registry.Alias("badPackage", &example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no package path")
	suite.Assert().Empty(suite.reg.aliases)
	err = suite.registry.Alias("x", example)
	suite.Assert().NoError(err)
	suite.Assert().Len(suite.reg.aliases, 1)
	err = suite.registry.Alias("x", example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "unable to redefine")
}

func (suite *RegistryTestSuite) TestRegister() {
	example := &alpha{}
	err := suite.registry.Register(&example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
	suite.Assert().Empty(suite.reg.byName)
	suite.Assert().Empty(suite.reg.byType)
	err = suite.registry.Register(example)
	suite.Assert().NoError(err)
	suite.Assert().Len(suite.reg.byName, 1)
	suite.Assert().Len(suite.reg.byType, 1)
	err = suite.registry.Register(example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "previous registration")
}

func (suite *RegistryTestSuite) TestNameFor() {
	example := &alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	exType, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(TypeUtilsPackagePath+"/alpha", exType)
}

func (suite *RegistryTestSuite) TestMake() {
	example := &alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	item, err := suite.registry.Make(TypeUtilsPackagePath + "/alpha")
	suite.Assert().NoError(err)
	suite.Assert().NotNil(item)
	suite.Assert().IsType(example, item)
}

func (suite *RegistryTestSuite) TestConverItemToMap() {
	suite.Assert().NoError(suite.registry.Register(&alpha{}))
	m, err := suite.registry.ConvertItemToMap(&alpha{
		Name:    "Goober Snoofus",
		Percent: 17.23,
		extra:   "nothing to see here",
	})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(m)
	suite.Assert().Len(m, 3)
	suite.Assert().Equal(TypeUtilsPackagePath+"/alpha", m[TypeField])
	suite.Assert().Equal("Goober Snoofus", m["Name"])
	suite.Assert().Equal(17.23, m["Percent"])
}

func (suite *RegistryTestSuite) TestCreateItemFromMap() {
	suite.Assert().NoError(suite.registry.Register(&alpha{}))
	example, err := suite.registry.CreateItemFromMap(map[string]interface{}{
		TypeField: TypeUtilsPackagePath + "/alpha",
		"Name":    "Goober Snoofus",
		"Percent": 17.23,
		"extra":   "nothing to see here",
	})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(example)
	suite.Assert().IsType(&alpha{}, example)
	suite.Assert().Equal(&alpha{
		Name:    "Goober Snoofus",
		Percent: 17.23,
	}, example)
}

func (suite *RegistryTestSuite) TestCycleSimple() {
	example := &alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	registration := suite.reg.byType[reflect.TypeOf(example).Elem()]
	suite.Assert().NotNil(registration)
	suite.Assert().Len(registration.allNames, 1)
	name, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(TypeUtilsPackagePath+"/alpha", name)
	object, err := suite.registry.Make(name)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(object)
	suite.Assert().Equal(reflect.TypeOf(example), reflect.TypeOf(object))
}

func (suite *RegistryTestSuite) TestCycleAlias() {
	example := &alpha{}
	suite.Assert().NoError(suite.registry.Alias("typeUtils", example))
	suite.Assert().NoError(suite.registry.Register(example))
	exType := reflect.TypeOf(example)
	if exType.Kind() == reflect.Ptr {
		exType = exType.Elem()
	}
	registration, ok := suite.reg.byType[exType]
	suite.Assert().True(ok)
	suite.Assert().NotNil(registration)

	suite.Assert().Len(registration.allNames, 2)
	name, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal("[typeUtils]alpha", name)
	object, err := suite.registry.Make(name)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(object)
	suite.Assert().Equal(reflect.TypeOf(example), reflect.TypeOf(object))
}

func (suite *RegistryTestSuite) TestGenTypeName() {
	example := &alpha{}
	name, err := genNameFromInterface(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(TypeUtilsPackagePath+"/alpha", name)
	_, err = genNameFromInterface(&example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
	_, err = genNameFromInterface(1)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
}
