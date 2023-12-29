package csv

import (
	baseCSV "encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/madkins23/go-utils/msg"
)

// Reader is a CSV reader based on encoding/csv.Reader.
// Additional functionality supports named access to fields.
type Reader struct {
	*baseCSV.Reader
	fieldLookup map[string]int
	fieldNames  []string
	indexLookup map[int]string
}

const errNoFieldNames msg.ConstError = "no field names provided"

// NewReader creates a new CSV reader object.
// At least one fieldName argument is required since the intent is to
// use them to access fields in records returned by Read.
// The wrapped encode/csv.Reader field is visible so its settings can be changed.
func NewReader(r io.Reader, fieldNames ...string) (*Reader, error) {
	if fieldNames == nil || len(fieldNames) < 1 {
		return nil, errNoFieldNames
	}

	rdr := &Reader{
		Reader:     baseCSV.NewReader(r),
		fieldNames: fieldNames,
	}
	rdr.Reader.Comment = '#'
	rdr.ReuseRecord = true

	fields, err := rdr.Reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read first line: %w", err)
	}

	// Track field names to column indexes.
	rdr.fieldLookup = make(map[string]int)
	for _, name := range rdr.fieldNames {
		rdr.fieldLookup[name] = -1
	}
	for i, v := range fields {
		if _, found := rdr.fieldLookup[v]; found {
			rdr.fieldLookup[v] = i
		}
	}

	// Make sure all known headers were found:
	missing := strings.Builder{}
	for name, i := range rdr.fieldLookup {
		if i < 0 {
			if missing.Len() > 0 {
				missing.WriteString(", ")
			}
			missing.WriteString(name)
		}
	}
	if missing.Len() > 0 {
		return nil, fmt.Errorf("first line missing headers: %s", missing.String())
	}

	// Create index lookup..
	rdr.indexLookup = make(map[int]string)
	for name, i := range rdr.fieldLookup {
		rdr.indexLookup[i] = name
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

// FieldName returns the name of the field specified by its column index.
// If the field is not named (not present in fieldNames in NewReader)
// an error will be return.
func (r *Reader) FieldName(index int) (string, error) {
	if name, found := r.indexLookup[index]; found {
		return name, nil
	} else {
		return "", fmt.Errorf("field index %d is not named", index)
	}
}

// Read consumes the next line and returns a map from field names to field values.
// Other errors may be returned as documented for encoding/csv.Reader.
// Specifically, at end of file the map returned is nil and the error is io.EOF.
func (r *Reader) Read() (map[string]string, error) {
	fields, err := r.Reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV line: %w", err)
	}

	// TODO: Add reuse record flag to this object.
	result := make(map[string]string)
	for i, v := range fields {
		if name, found := r.indexLookup[i]; found {
			result[name] = v
		}
	}
	return result, nil
}
