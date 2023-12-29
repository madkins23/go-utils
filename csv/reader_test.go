package csv

import (
	"encoding/csv"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	header  = `"alpha","dummy","bravo",,"charlie"`
	csvBody = `"one","two","three","four","five"
"1","2","3","4","5"`
	withHeader          = header + "\n" + csvBody
	tooFewFields        = `"two","few"`
	tooManyFields       = `"two","many","freaking","fields","in","the","row"`
	disorderedFirstLine = `"goober","bravo","charlie","snoofus","alpha"`
	oneMatchFirstLine   = `"goober","bravo","snoofus","dummy","names"`
	rowLength           = 3
)

var fieldNames = []string{
	"alpha",
	"bravo",
	"charlie",
}

func TestNewReader(t *testing.T) {
	makeReader(t, withHeader, fieldNames...)
}

func TestNewReader_noHeader(t *testing.T) {
	_, err := NewReader(strings.NewReader(csvBody))
	assert.ErrorIs(t, err, errNoFieldNames)
	_, err = NewReader(strings.NewReader(csvBody), []string{}...)
	assert.ErrorIs(t, err, errNoFieldNames)
}

func TestReader_FieldIndex(t *testing.T) {
	reader := makeReader(t, withHeader, fieldNames...)
	index, err := reader.FieldIndex("alpha")
	require.NoError(t, err)
	assert.Equal(t, 0, index)
	index, err = reader.FieldIndex("bravo")
	require.NoError(t, err)
	assert.Equal(t, 2, index)
	index, err = reader.FieldIndex("charlie")
	require.NoError(t, err)
	assert.Equal(t, 4, index)
}

func TestReader_FieldIndex_badName(t *testing.T) {
	reader := makeReader(t, withHeader, fieldNames...)
	index, err := reader.FieldIndex("")
	assert.ErrorContains(t, err, "field name '' not found")
	assert.Equal(t, -1, index)
	index, err = reader.FieldIndex("goober")
	assert.ErrorContains(t, err, "field name 'goober' not found")
	assert.Equal(t, -1, index)
}

func TestReader_FieldName(t *testing.T) {
	reader := makeReader(t, withHeader, fieldNames...)
	index, err := reader.FieldName(0)
	require.NoError(t, err)
	assert.Equal(t, "alpha", index)
	index, err = reader.FieldName(2)
	require.NoError(t, err)
	assert.Equal(t, "bravo", index)
	index, err = reader.FieldName(4)
	require.NoError(t, err)
	assert.Equal(t, "charlie", index)
}

func TestReader_FieldName_badIndex(t *testing.T) {
	reader := makeReader(t, withHeader, fieldNames...)
	name, err := reader.FieldName(-1)
	assert.ErrorContains(t, err, "field index -1 is not named")
	assert.Equal(t, "", name)
	name, err = reader.FieldName(3)
	assert.ErrorContains(t, err, "field index 3 is not named")
	assert.Equal(t, "", name)
	name, err = reader.FieldName(9)
	assert.ErrorContains(t, err, "field index 9 is not named")
	assert.Equal(t, "", name)
}

func TestReader(t *testing.T) {
	reader := makeReader(t, withHeader, fieldNames...)
	row, err := reader.Read()
	assert.NoError(t, err)
	assert.Len(t, row, rowLength)
	assert.Equal(t, "one", row["alpha"])
	assert.Equal(t, "three", row["bravo"])
	assert.Equal(t, "five", row["charlie"])
	row, err = reader.Read()
	assert.NoError(t, err)
	assert.Len(t, row, rowLength)
	assert.Equal(t, "1", row["alpha"])
	assert.Equal(t, "3", row["bravo"])
	assert.Equal(t, "5", row["charlie"])
	row, err = reader.Read()
	assert.ErrorContains(t, err, "read CSV line")
	assert.ErrorAs(t, err, &io.EOF)
	assert.Nil(t, row)
}

func TestReader_Read_justHeader(t *testing.T) {
	reader := makeReader(t, header, fieldNames...)
	row, err := reader.Read()
	assert.ErrorIs(t, err, io.EOF)
	assert.Nil(t, row)
}

func TestReader_Read_tooFewFields(t *testing.T) {
	reader := makeReader(t, header+"\n"+tooFewFields, fieldNames...)
	row, err := reader.Read()
	assert.ErrorIs(t, err, csv.ErrFieldCount)
	assert.Nil(t, row)
}

func TestReader_Read_tooManyFields(t *testing.T) {
	reader := makeReader(t, header+"\n"+tooManyFields, fieldNames...)
	row, err := reader.Read()
	assert.ErrorContains(t, err, "wrong number of fields")
	assert.Nil(t, row)
}

func TestReader_Read_noHeader(t *testing.T) {
	reader, err := NewReader(strings.NewReader(csvBody), fieldNames...)
	assert.ErrorContains(t, err, "first line missing headers")
	assert.ErrorContains(t, err, "alpha")
	assert.ErrorContains(t, err, "bravo")
	assert.ErrorContains(t, err, "charlie")
	assert.Nil(t, reader)
}

func makeReader(t *testing.T, csvData string, fieldNames ...string) *Reader {
	reader, err := NewReader(strings.NewReader(csvData), fieldNames...)
	require.NoError(t, err)
	require.NotNil(t, reader)
	return reader
}
