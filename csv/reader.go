package csv

import (
	baseCSV "encoding/csv"
	"fmt"
	"io"

	"github.com/madkins23/go-utils/msg"
)

// Reader is a CSV reader based on encoding/csv.Reader.
// Additional functionality supports named access to fields.
type Reader struct {
	*baseCSV.Reader
	firstLine   bool
	fieldNames  []string
	fieldLookup map[string]int
	numFields   int
}

const errNoFieldNames msg.ConstError = "no field names provided"

// NewReader creates a new CSV reader object.
// The fieldNames argument is required and must not be nil or empty.
// If no field names are available it is better to use encoding/csv.Reader directly.
// The wrapped encode/csv.Reader field is visible so its settings can be changed.
func NewReader(r io.Reader, fieldNames []string) (*Reader, error) {
	if fieldNames == nil || len(fieldNames) < 1 {
		return nil, errNoFieldNames
	}
	numFields := len(fieldNames)

	rdr := &Reader{
		Reader:      baseCSV.NewReader(r),
		firstLine:   true,
		fieldNames:  fieldNames,
		fieldLookup: make(map[string]int),
		numFields:   numFields,
	}
	rdr.FieldsPerRecord = numFields
	rdr.Reader.Comment = '#'
	rdr.ReuseRecord = true
	for i, v := range fieldNames {
		rdr.fieldLookup[v] = i
	}

	return rdr, nil
}

// FieldIndex returns the index of the field specified by name.
// An error is returned if there is no such field name.
func (r *Reader) FieldIndex(name string) (int, error) {
	if index, found := r.fieldLookup[name]; !found {
		return -1, fmt.Errorf("field name '%s' not found", name)
	} else {
		return index, nil
	}
}

// FieldName returns the name of the field specified by index.
// An error is returned if the field is out of range.
func (r *Reader) FieldName(index int) (string, error) {
	if index < 0 || index >= r.numFields {
		return "", fmt.Errorf("field index %d out of range", index)
	} else {
		return r.fieldNames[index], nil
	}
}

// Read consumes the next line and returns a map from field names to field values.
// Other errors may be returned as documented for encoding/csv.Reader.
// Specifically, at end of file the map returned is nil and the error is io.EOF.
func (r *Reader) Read() (record map[string]string, err error) {
	if fields, err := r.Reader.Read(); err != nil {
		return nil, err
	} else {
		if r.firstLine {
			// Attempt to check first line to see if it has correct column names.
			var matches, disordered int
			for i, v := range fields {
				if v == r.fieldNames[i] {
					matches++
				} else if _, found := r.fieldLookup[v]; found {
					disordered++
				}
			}
			if matches > 0 || disordered > 0 {
				// Found at least one column name, check further.
				if matches != r.numFields || disordered > 0 {
					return nil, fmt.Errorf("first line error: %d matches, %d disordered", matches, disordered)
				}
				// This was column header line so read another one for first data line.
				if fields, err = r.Reader.Read(); err != nil {
					return nil, err
				}
			}
			r.firstLine = false
		}

		// TODO: Add reuse record flag to this object.
		result := make(map[string]string)
		for i, v := range fields {
			result[r.fieldNames[i]] = v
		}
		return result, nil
	}
}
