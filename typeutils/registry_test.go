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
	suite.registry = NewRegistry()
	var ok bool
	suite.reg, ok = suite.registry.(*registry)
	suite.Assert().True(ok)
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

func (suite *RegistryTestSuite) TestSimpleCases() {
	example := alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	registration := suite.reg.byType[reflect.TypeOf(example)]
	suite.Assert().Len(registration.allNames, 1)
	name, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(PKG_PATH+"/alpha", name)
	object, err := suite.registry.Make(name)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(object)
	suite.Assert().Equal(reflect.TypeOf(example), reflect.TypeOf(object))
}

func (suite *RegistryTestSuite) TestAliasCases() {
	example := alpha{}
	suite.Assert().NoError(suite.registry.Alias("test", example))
	suite.Assert().NoError(suite.registry.Register(example))
	registration := suite.reg.byType[reflect.TypeOf(example)]
	suite.Assert().Len(registration.allNames, 2)
	name, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal("[test]alpha", name)
	object, err := suite.registry.Make(name)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(object)
	suite.Assert().Equal(reflect.TypeOf(example), reflect.TypeOf(object))
}
