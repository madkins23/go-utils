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
	header   = `"alpha","bravo","charlie"`
	noHeader = `"one","two","three"
"1","2","3"`
	tooFewFields        = `"two","few"`
	tooManyFields       = `"two","many","freaking","fields"`
	disorderedFirstLine = `"goober","alpha","snoofus"`
	oneMatchFirstLine   = `"goober","bravo","snoofus"`
)

var fieldNames []string

func init() {
	names := strings.Split(header, ",")
	fieldNames = make([]string, len(names))
	for i, name := range names {
		fieldNames[i] = strings.Trim(name, "\"")
	}
}

func TestNewReader(t *testing.T) {
	makeReader(t, noHeader, fieldNames)
}

func TestNewReader_noHeader(t *testing.T) {
	_, err := NewReader(strings.NewReader(noHeader), nil)
	assert.ErrorIs(t, err, errNoFieldNames)
	_, err = NewReader(strings.NewReader(noHeader), []string{})
	assert.ErrorIs(t, err, errNoFieldNames)
}

func TestReader_FieldIndex(t *testing.T) {
	reader := makeReader(t, noHeader, fieldNames)
	index, err := reader.FieldIndex("alpha")
	require.NoError(t, err)
	assert.Equal(t, 0, index)
	index, err = reader.FieldIndex("bravo")
	require.NoError(t, err)
	assert.Equal(t, 1, index)
	index, err = reader.FieldIndex("charlie")
	require.NoError(t, err)
	assert.Equal(t, 2, index)
}

func TestReader_FieldIndex_badName(t *testing.T) {
	reader := makeReader(t, noHeader, fieldNames)
	index, err := reader.FieldIndex("")
	assert.ErrorContains(t, err, "field name '' not found")
	assert.Equal(t, -1, index)
	index, err = reader.FieldIndex("goober")
	assert.ErrorContains(t, err, "field name 'goober' not found")
	assert.Equal(t, -1, index)
}

func TestReader_FieldName(t *testing.T) {
	reader := makeReader(t, noHeader, fieldNames)
	index, err := reader.FieldName(0)
	require.NoError(t, err)
	assert.Equal(t, "alpha", index)
	index, err = reader.FieldName(1)
	require.NoError(t, err)
	assert.Equal(t, "bravo", index)
	index, err = reader.FieldName(2)
	require.NoError(t, err)
	assert.Equal(t, "charlie", index)
}

func TestReader_FieldName_badIndex(t *testing.T) {
	reader := makeReader(t, noHeader, fieldNames)
	name, err := reader.FieldName(-1)
	assert.ErrorContains(t, err, "field index -1 out of range")
	assert.Equal(t, "", name)
	name, err = reader.FieldName(3)
	assert.ErrorContains(t, err, "field index 3 out of range")
	assert.Equal(t, "", name)
}

func TestReader_Read_justHeader(t *testing.T) {
	reader := makeReader(t, header, fieldNames)
	row, err := reader.Read()
	assert.ErrorIs(t, err, io.EOF)
	assert.Nil(t, row)
}

func TestReader_Read_tooFewFields(t *testing.T) {
	reader := makeReader(t, tooFewFields, fieldNames)
	row, err := reader.Read()
	assert.ErrorIs(t, err, csv.ErrFieldCount)
	assert.Nil(t, row)
}

func TestReader_Read_tooManyFields(t *testing.T) {
	reader := makeReader(t, tooManyFields, fieldNames)
	row, err := reader.Read()
	assert.ErrorIs(t, err, csv.ErrFieldCount)
	assert.Nil(t, row)
}

func TestReader_Read_firstLineDisordered(t *testing.T) {
	reader := makeReader(t, disorderedFirstLine, fieldNames)
	_, err := reader.Read()
	require.ErrorContains(t, err, "0 matches, 1 disordered")
}

func TestReader_Read_firstLineOneMatch(t *testing.T) {
	reader := makeReader(t, oneMatchFirstLine, fieldNames)
	_, err := reader.Read()
	require.ErrorContains(t, err, "1 matches, 0 disordered")
}

func TestReader_Read_noHeader(t *testing.T) {
	reader := makeReader(t, header+"\n"+noHeader, fieldNames)
	readNoHeader(t, reader)
}

func TestReader_Read_withHeader(t *testing.T) {
	reader := makeReader(t, noHeader, fieldNames)
	readNoHeader(t, reader)
}

func readNoHeader(t *testing.T, reader *Reader) {
	row, err := reader.Read()
	require.NoError(t, err)
	assert.Equal(t, "one", row["alpha"])
	assert.Equal(t, "two", row["bravo"])
	assert.Equal(t, "three", row["charlie"])
	row, err = reader.Read()
	require.NoError(t, err)
	assert.Equal(t, "1", row["alpha"])
	assert.Equal(t, "2", row["bravo"])
	assert.Equal(t, "3", row["charlie"])
	_, err = reader.Read()
	require.Error(t, err, io.EOF)
}

func makeReader(t *testing.T, csvData string, fieldNames []string) *Reader {
	reader, err := NewReader(strings.NewReader(csvData), fieldNames)
	require.NoError(t, err)
	require.NotNil(t, reader)
	return reader
}

func withHeader(body string) string {
	return header + "\n" + body
}
