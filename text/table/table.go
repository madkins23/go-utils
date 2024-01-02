package table

import (
	"fmt"
	"strings"
)

// Find Unicode characters using:
//  https://en.wikipedia.org/wiki/Box-drawing_character
//  https://www.fileformat.info/info/unicode/char/253C/index.htm

const (
	newLine byte = '\n'
	space   byte = ' '
)

type axis byte

const (
	horizontal axis = iota
	vertical
)

type horizontals byte

const (
	left horizontals = iota
	right
)

type Verticals byte

const (
	Top Verticals = iota
	Bottom
)

// NumLines provides an enumerated type to specify single or double lines.
type NumLines byte

const (
	Single NumLines = iota
	Double
)

var (
	// Runes for corners: [Verticals][horizontals][NumLines]
	corner [][][]rune

	// Runes for line crossings within table:  [direction][NumLines]
	cross [][]rune

	// Runes for line segments:   [horizontal NumLines][vertical NumLines]
	line [][]rune

	// Runes for lines ending in lines (i.e. like a 'T'): [axis][Verticals][bar NumLines][leg NumLines]
	// The horizontal top line of a 'T' is the "bar" and the vertical line is the "leg".
	tee [][][][]rune
)

func init() {
	corner = make([][][]rune, 2)
	corner[Top] = make([][]rune, 2)
	corner[Top][left] = []rune{'┌', '╔'}
	corner[Top][right] = []rune{'┐', '╗'}
	corner[Bottom] = make([][]rune, 2)
	corner[Bottom][left] = []rune{'└', '╚'}
	corner[Bottom][right] = []rune{'┘', '╝'}

	cross = make([][]rune, 2)
	cross[Single] = []rune{'┼', '╫'}
	cross[Double] = []rune{'╪', '╬'}

	line = make([][]rune, 2)
	line[horizontal] = []rune{'─', '═'}
	line[vertical] = []rune{'│', '║'}

	tee = make([][][][]rune, 4)
	tee[horizontal] = make([][][]rune, 2)
	tee[horizontal][left] = [][]rune{{'├', '╞'}, {'╟', '╠'}}
	tee[horizontal][right] = [][]rune{{'┤', '╡'}, {'╢', '╣'}}
	tee[vertical] = make([][][]rune, 2)
	tee[vertical][Top] = [][]rune{{'┬', '╥'}, {'╤', '╦'}}
	tee[vertical][Bottom] = [][]rune{{'┴', '╨'}, {'╧', '╩'}}
}

// ColumnDef defines a column for the table.
type ColumnDef struct {
	// The width of the column in spaces.
	Width int

	// The printf format specification for row data for the column.
	// If this is left empty a simple string format (e.g. '%4s') will be used.
	Format string

	ColumnLines NumLines

	// Specify alignment of column header and row contents if Format is not specified.
	// By default, strings alight to the right.
	AlignLeft bool
}

// TableDef defines a table.
type TableDef struct {
	// Array of ColumnDef objects that specify the columns for the table.
	Columns []ColumnDef

	// Prefix string for each row.
	Prefix string

	// Border presence
	Border bool

	// Border lines
	BorderLines NumLines

	header, row       string
	border, separator []string
}

// BorderString returns a string representing a Top border row.
// This should only be used if the specified TableDef has its Border flag set.
// No terminal newline is attached as this string is intended to be used with println.
func (td *TableDef) BorderString(v Verticals) string {
	if td.border == nil {
		td.border = make([]string, 2)
	}

	if td.border[v] == "" {
		builder := td.newBuilder()
		var started bool

		builder.WriteRune(corner[v][left][td.BorderLines])
		builder.WriteRune(line[horizontal][td.BorderLines])

		for _, column := range td.Columns {
			if started {
				builder.WriteRune(line[horizontal][td.BorderLines])
				builder.WriteRune(tee[vertical][v][td.BorderLines][column.ColumnLines])
				builder.WriteRune(line[horizontal][td.BorderLines])
			} else {
				started = true
			}
			for i := 0; i < column.Width; i++ {
				builder.WriteRune(line[horizontal][td.BorderLines])
			}
		}

		builder.WriteRune(line[horizontal][td.BorderLines])
		builder.WriteRune(corner[v][right][td.BorderLines])

		td.border[v] = builder.String()
	}

	return td.border[v]
}

// HeaderFormat returns a printf format string for the header row.
// This string includes a terminal newline.
func (td *TableDef) HeaderFormat() string {
	if td.header == "" {
		builder := td.newBuilder()
		var started bool

		if td.Border {
			builder.WriteRune(line[vertical][td.BorderLines])
			builder.WriteByte(space)
		}

		for _, column := range td.Columns {
			if started {
				builder.WriteByte(space)
				builder.WriteRune(line[vertical][column.ColumnLines])
				builder.WriteByte(space)
			} else {
				started = true
			}
			builder.WriteString(column.headerFormat())
		}

		if td.Border {
			builder.WriteByte(space)
			builder.WriteRune(line[vertical][td.BorderLines])
		}

		builder.WriteByte(newLine)
		td.header = builder.String()
	}

	return td.header
}

// DividerString returns a string representing a separator row.
// No terminal newline is attached as this string is intended to be used with println.
// Deprecated: Use SeparatorString instead
func (td *TableDef) DividerString() string {
	return td.SeparatorString(Single)
}

// SeparatorString returns a string representing a separator row with the specified number of lines.
// No terminal newline is attached as this string is intended to be used with println.
func (td *TableDef) SeparatorString(lines NumLines) string {
	if td.separator == nil {
		td.separator = make([]string, 2)
	}

	if td.separator[lines] == "" {
		builder := td.newBuilder()
		var started bool

		if td.Border {
			builder.WriteRune(tee[horizontal][left][td.BorderLines][lines])
			builder.WriteRune(line[horizontal][lines])
		}

		for _, column := range td.Columns {
			if started {
				builder.WriteRune(line[horizontal][lines])
				builder.WriteRune(cross[lines][column.ColumnLines])
				builder.WriteRune(line[horizontal][lines])
			} else {
				started = true
			}
			for i := 0; i < column.Width; i++ {
				builder.WriteRune(line[horizontal][lines])
			}
		}

		if td.Border {
			builder.WriteRune(line[horizontal][lines])
			builder.WriteRune(tee[horizontal][right][td.BorderLines][lines])
		}

		td.separator[lines] = builder.String()
	}

	return td.separator[lines]
}

// RowFormat returns a printf format string for a data row.
// This string includes a terminal newline.
func (td *TableDef) RowFormat() string {
	if td.row == "" {
		builder := td.newBuilder()
		var started bool

		if td.Border {
			builder.WriteRune(line[vertical][td.BorderLines])
			builder.WriteByte(space)
		}

		for _, column := range td.Columns {
			if started {
				builder.WriteByte(space)
				builder.WriteRune(line[vertical][column.ColumnLines])
				builder.WriteByte(space)
			} else {
				started = true
			}
			builder.WriteString(column.dataFormat())
		}

		if td.Border {
			builder.WriteByte(space)
			builder.WriteRune(line[vertical][td.BorderLines])
		}

		builder.WriteByte(newLine)
		td.row = builder.String()
	}

	return td.row
}

func (td *TableDef) newBuilder() *strings.Builder {
	var builder strings.Builder

	if td.Prefix != "" {
		builder.WriteString(td.Prefix)
	}

	return &builder
}

func (c *ColumnDef) dataFormat() string {
	if c.Format != "" {
		return c.Format
	} else {
		return c.headerFormat()
	}
}

func (c *ColumnDef) headerFormat() string {
	if c.AlignLeft {
		return fmt.Sprintf("%%-%ds", c.Width)
	} else {
		return fmt.Sprintf("%%%ds", c.Width)
	}
}
