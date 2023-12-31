package text

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TableTestSuite struct {
	suite.Suite
	tableDef *TableDef
}

func (suite *TableTestSuite) SetupTest() {
	suite.tableDef = &TableDef{
		Columns: []ColumnDef{
			{Width: 5, AlignLeft: true},
			{Width: 10, Format: "%10d"},
			{Width: 8, Format: "%8f", Double: true},
		},
		Prefix: "> ",
	}
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TableTestSuite))
}

func (suite *TableTestSuite) TestTable_Header() {
	assert.Equal(suite.T(), "> %-5s │ %10s ║ %8s\n", suite.tableDef.HeaderFormat())
}

func (suite *TableTestSuite) TestTable_Divider() {
	assert.Equal(suite.T(), "> ──────┼────────────╫─────────", suite.tableDef.DividerString())
}

func (suite *TableTestSuite) TestTable_Row() {
	assert.Equal(suite.T(), "> %-5s │ %10d ║ %8f\n", suite.tableDef.RowFormat())
}

func ExampleTableDef() {
	tableDef := &TableDef{
		Columns: []ColumnDef{
			{Width: 5, AlignLeft: true},
			{Width: 10, Format: "%10d"},
			{Width: 8, Format: "%8f", Double: true},
		},
		Prefix: "> ",
	}

	fmt.Printf(tableDef.HeaderFormat(), "name", "count", "float")
	fmt.Println(tableDef.DividerString())
	fmt.Printf(tableDef.RowFormat(), "x", 3, 4.5)

	// Output: > name  │      count ║    float
	//> ──────┼────────────╫─────────
	//> x     │          3 ║ 4.500000
}
