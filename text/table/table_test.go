package table

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
			{Width: 8, Format: "%8f", ColumnLines: Double},
		},
		Prefix: "> ",
	}
}

func TestTableSuite(t *testing.T) {
	suite.Run(t, new(TableTestSuite))
}

func (suite *TableTestSuite) TestTable_BorderSingle() {
	assert.Equal(suite.T(), "> ┌───────┬────────────╥──────────┐", suite.tableDef.BorderString(Top))
	assert.Equal(suite.T(), "> └───────┴────────────╨──────────┘", suite.tableDef.BorderString(Bottom))
}

func (suite *TableTestSuite) TestTable_BorderDouble() {
	suite.tableDef.BorderLines = Double
	assert.Equal(suite.T(), "> ╔═══════╤════════════╦══════════╗", suite.tableDef.BorderString(Top))
	assert.Equal(suite.T(), "> ╚═══════╧════════════╩══════════╝", suite.tableDef.BorderString(Bottom))
}

func (suite *TableTestSuite) TestTable_Header() {
	assert.Equal(suite.T(), "> %-5s │ %10s ║ %8s\n", suite.tableDef.HeaderFormat())
}

func (suite *TableTestSuite) TestTable_Header_borderSingle() {
	suite.tableDef.Border = true
	assert.Equal(suite.T(), "> │ %-5s │ %10s ║ %8s │\n", suite.tableDef.HeaderFormat())
}

func (suite *TableTestSuite) TestTable_Header_borderDouble() {
	suite.tableDef.Border = true
	suite.tableDef.BorderLines = Double
	assert.Equal(suite.T(), "> ║ %-5s │ %10s ║ %8s ║\n", suite.tableDef.HeaderFormat())
}

func (suite *TableTestSuite) TestTable_Separator() {
	assert.Equal(suite.T(), "> ──────┼────────────╫─────────", suite.tableDef.SeparatorString(Single))
	assert.Equal(suite.T(), "> ══════╪════════════╬═════════", suite.tableDef.SeparatorString(Double))
}

func (suite *TableTestSuite) TestTable_Separator_borderSingle() {
	suite.tableDef.Border = true
	assert.Equal(suite.T(), "> ├───────┼────────────╫──────────┤", suite.tableDef.SeparatorString(Single))
	assert.Equal(suite.T(), "> ╞═══════╪════════════╬══════════╡", suite.tableDef.SeparatorString(Double))
}

func (suite *TableTestSuite) TestTable_Separator_borderDouble() {
	suite.tableDef.Border = true
	suite.tableDef.BorderLines = Double
	assert.Equal(suite.T(), "> ╟───────┼────────────╫──────────╢", suite.tableDef.SeparatorString(Single))
	assert.Equal(suite.T(), "> ╠═══════╪════════════╬══════════╣", suite.tableDef.SeparatorString(Double))
}

func (suite *TableTestSuite) TestTable_Row() {
	assert.Equal(suite.T(), "> %-5s │ %10d ║ %8f\n", suite.tableDef.RowFormat())
}

func (suite *TableTestSuite) TestTable_Row_borderSingle() {
	suite.tableDef.Border = true
	assert.Equal(suite.T(), "> │ %-5s │ %10d ║ %8f │\n", suite.tableDef.RowFormat())
}

func (suite *TableTestSuite) TestTable_Row_borderDouble() {
	suite.tableDef.Border = true
	suite.tableDef.BorderLines = Double
	assert.Equal(suite.T(), "> ║ %-5s │ %10d ║ %8f ║\n", suite.tableDef.RowFormat())
}

func ExampleTableDef() {
	table := &TableDef{
		Columns: []ColumnDef{
			{Width: 5, AlignLeft: true},
			{Width: 10, Format: "%10d"},
			{Width: 8, Format: "%8.2f", ColumnLines: Double},
		},
		Prefix: "> ",
	}

	fmt.Println()
	fmt.Printf(table.HeaderFormat(), "name", "count", "float")
	fmt.Println(table.SeparatorString(Single))
	fmt.Printf(table.RowFormat(), "x", 3, 4.5)
	fmt.Println(table.SeparatorString(Double))
	fmt.Printf(table.RowFormat(), "z", 8, 9.2)

	// Output:
	//> name  │      count ║    float
	//> ──────┼────────────╫─────────
	//> x     │          3 ║     4.50
	//> ══════╪════════════╬═════════
	//> z     │          8 ║     9.20
}

func ExampleTableDef_borderSingle() {
	table := &TableDef{
		Columns: []ColumnDef{
			{Width: 5, AlignLeft: true},
			{Width: 10, Format: "%10d"},
			{Width: 8, Format: "%8.2f", ColumnLines: Double},
		},
		Prefix: "> ",
		Border: true,
	}

	fmt.Println()
	fmt.Println(table.BorderString(Top))
	fmt.Printf(table.HeaderFormat(), "name", "count", "float")
	fmt.Println(table.SeparatorString(Single))
	fmt.Printf(table.RowFormat(), "x", 3, 4.5)
	fmt.Println(table.SeparatorString(Double))
	fmt.Printf(table.RowFormat(), "z", 8, 9.2)
	fmt.Println(table.BorderString(Bottom))

	// Output:
	//> ┌───────┬────────────╥──────────┐
	//> │ name  │      count ║    float │
	//> ├───────┼────────────╫──────────┤
	//> │ x     │          3 ║     4.50 │
	//> ╞═══════╪════════════╬══════════╡
	//> │ z     │          8 ║     9.20 │
	//> └───────┴────────────╨──────────┘
}

func ExampleTableDef_borderDouble() {
	table := &TableDef{
		Columns: []ColumnDef{
			{Width: 5, AlignLeft: true},
			{Width: 10, Format: "%10d"},
			{Width: 8, Format: "%8.2f", ColumnLines: Double},
		},
		Prefix:      "> ",
		Border:      true,
		BorderLines: Double,
	}

	fmt.Println()
	fmt.Println(table.BorderString(Top))
	fmt.Printf(table.HeaderFormat(), "name", "count", "float")
	fmt.Println(table.SeparatorString(Single))
	fmt.Printf(table.RowFormat(), "x", 3, 4.5)
	fmt.Println(table.SeparatorString(Double))
	fmt.Printf(table.RowFormat(), "z", 8, 9.2)
	fmt.Println(table.BorderString(Bottom))

	// Output:
	//> ╔═══════╤════════════╦══════════╗
	//> ║ name  │      count ║    float ║
	//> ╟───────┼────────────╫──────────╢
	//> ║ x     │          3 ║     4.50 ║
	//> ╠═══════╪════════════╬══════════╣
	//> ║ z     │          8 ║     9.20 ║
	//> ╚═══════╧════════════╩══════════╝
}
