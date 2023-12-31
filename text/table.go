package text

import (
	"fmt"
	"strings"
)

// Find Unicode characters using:
//  https://en.wikipedia.org/wiki/Box-drawing_character
//  https://www.fileformat.info/info/unicode/char/253C/index.htm

const (
	doubleVertical                 rune = '\u2551' // '\U00E28FB8'
	singleHorizontal                    = '\u2500'
	singleHorizontalDoubleVertical      = '\u256B'
	singleHorizontalSingleVertical      = '\u253C'
	singleVertical                      = '\u2502'

	newLine byte = '\n'
	space   byte = ' '
)

// ColumnDef defines a column for the table.
type ColumnDef struct {
	// The width of the column in spaces.
	Width int

	// The printf format specification for row data for the column.
	// If this is left empty a simple string format (e.g. '%4s') will be used.
	Format string

	// For columns after the first specifies that the vertical separator
	// character before the column should be double rather than single.
	Double bool

	// Specify alignment of column header and row contents if Format is not specified.
	// By default strings alight to the right.
	AlignLeft bool
}

// TableDef defines a table.
type TableDef struct {
	// Array of ColumnDef objects that specify the columns for the table.
	Columns []ColumnDef

	// Prefix string for each row.
	Prefix string

	header, divider, row string
}

// HeaderFormat returns a printf format string for the header row.
// This string includes a terminal newline.
func (td *TableDef) HeaderFormat() string {
	if td.header == "" {
		builder := td.newBuilder()
		var started bool

		for _, column := range td.Columns {
			if started {
				builder.WriteByte(space)
				if column.Double {
					builder.WriteRune(doubleVertical)
				} else {
					builder.WriteRune(singleVertical)
				}
				builder.WriteByte(space)
			} else {
				started = true
			}
			builder.WriteString(column.headerFormat())
		}

		builder.WriteByte(newLine)
		td.header = builder.String()
	}

	return td.header
}

// DividerString returns a string representing a divider row.
// No terminal newline is attached as this string is intended to be used with println.
func (td *TableDef) DividerString() string {
	if td.divider == "" {
		builder := td.newBuilder()
		var started bool

		for _, column := range td.Columns {
			if started {
				builder.WriteRune(singleHorizontal)
				if column.Double {
					builder.WriteRune(singleHorizontalDoubleVertical)
				} else {
					builder.WriteRune(singleHorizontalSingleVertical)
				}
				builder.WriteRune(singleHorizontal)
			} else {
				started = true
			}
			for i := 0; i < column.Width; i++ {
				builder.WriteRune(singleHorizontal)
			}
		}

		td.divider = builder.String()
	}

	return td.divider
}

// RowFormat returns a printf format string for a data row.
// This string includes a terminal newline.
func (td *TableDef) RowFormat() string {
	if td.row == "" {
		builder := td.newBuilder()
		var started bool

		for _, column := range td.Columns {
			if started {
				builder.WriteByte(space)
				if column.Double {
					builder.WriteRune(doubleVertical)
				} else {
					builder.WriteRune(singleVertical)
				}
				builder.WriteByte(space)
			} else {
				started = true
			}
			builder.WriteString(column.dataFormat())
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
