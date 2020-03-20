package typeutils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	a        = alpha{what: "dunno"}
	b        = bravo{ok: true}
	c        = a
	ai Actor = &a
	bi Actor = &b
	ci Actor = &c
)

//////////////////////////////////////////////////////////////////////////

type ReflectTestSuite struct {
	suite.Suite
}

func (suite *ReflectTestSuite) SetupTest() {
}
func TestReflectSuite(t *testing.T) {
	suite.Run(t, new(ReflectTestSuite))
}

//////////////////////////////////////////////////////////////////////////
// Verify method for determining path of package via an object defined therein.

// PKG_PATH should be set to the known package path for this package.
const PKG_PATH = "github.com/madkins23/go-utils/typeutils"

func (suite *ReflectTestSuite) TestPackagePath() {
	suite.Assert().Equal(PKG_PATH, reflect.TypeOf(alpha{}).PkgPath())
}

//////////////////////////////////////////////////////////////////////////
// Make certain reflect.Type supports equivalence testing and use as map key.
// Note that this does NOT work for types.Type, which is a different thing.

func (suite *ReflectTestSuite) TestTypeEquivalence() {
	suite.Assert().NotEqual(a, b)
	suite.Assert().NotEqual(b, c)
	suite.Assert().Equal(a, c)
}

func (suite *ReflectTestSuite) TestInterfaceEquivalence() {
	suite.Assert().NotEqual(ai, bi)
	suite.Assert().NotEqual(bi, ci)
	suite.Assert().Equal(ai, ci)
}

func (suite *ReflectTestSuite) TestMap() {
	lookup := make(map[reflect.Type]string)
	lookup[reflect.TypeOf(a)] = "alpha"
	lookup[reflect.TypeOf(b)] = "bravo"
	lookup[reflect.TypeOf(c)] = "charlie" // overrides "alpha" since a == c
	suite.Assert().Equal("charlie", lookup[reflect.TypeOf(a)])
	suite.Assert().Equal("bravo", lookup[reflect.TypeOf(b)])
	suite.Assert().Equal("charlie", lookup[reflect.TypeOf(c)])
}

func (suite *ReflectTestSuite) TestMapInterface() {
	lookup := make(map[reflect.Type]string)
	lookup[reflect.TypeOf(ai)] = "alpha"
	lookup[reflect.TypeOf(bi)] = "bravo"
	lookup[reflect.TypeOf(ci)] = "charlie" // overrides "alpha" since ai == ci
	suite.Assert().Equal("charlie", lookup[reflect.TypeOf(ai)])
	suite.Assert().Equal("bravo", lookup[reflect.TypeOf(bi)])
	suite.Assert().Equal("charlie", lookup[reflect.TypeOf(ci)])
}
