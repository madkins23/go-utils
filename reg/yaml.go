package reg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type YamlBase struct {
	Registry
}

func NewYamlBase(reg Registry) *YamlBase {
	return &YamlBase{Registry: reg}
}

//////////////////////////////////////////////////////////////////////////

// LoadFromFile loads an item of a registered type from the specified YAML file.
func (yb *YamlBase) LoadFromFile(fileName string) (item interface{}, finalErr error) {
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open source file %s: %w", fileName, err)
	}
	defer func() {
		if err = reader.Close(); err != nil {
			finalErr = fmt.Errorf("close file: %w", err)
		}
	}()

	if item, err = yb.loadFromReadSeeker(reader); finalErr != nil {
		return nil, fmt.Errorf("load from readSeeker: %w", err)
	} else {
		return item, nil
	}
}

// LoadFromString loads an item of a registered type from a YAML string.
func (yb *YamlBase) LoadFromString(source string) (interface{}, error) {
	if nexus, err := yb.loadFromReadSeeker(strings.NewReader(source)); err != nil {
		return nil, fmt.Errorf("load from readSeeker: %w", err)
	} else {
		return nexus, nil
	}
}

// loadFromReadSeeker loads a registered type from an open io.ReadSeeker.
// This could have been an io.Reader but it is necessary to 'parse' the data twice.
// The first time is just a scan for the top-level object type.
// After scanning partially through the stream it must be reset.
func (yb *YamlBase) loadFromReadSeeker(reader io.ReadSeeker) (interface{}, error) {
	if typeName, err := getYamlTypeNameAndReset(reader); err != nil {
		return nil, fmt.Errorf("get YAML type name: %w", err)
	} else if item, err := yb.Registry.Make(typeName); err != nil {
		return nil, fmt.Errorf("make item of type %s: %w", typeName, err)
	} else if err := yaml.NewDecoder(reader).Decode(item); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", typeName, err)
	} else {
		return item, nil
	}
}

//////////////////////////////////////////////////////////////////////////

// SaveToFile saves an item of a registered type to the specified YAML file.
func (yb *YamlBase) SaveToFile(item interface{}, fileName string) (finalErr error) {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("creating output file '%s': %w", fileName, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			finalErr = fmt.Errorf("Error closing source file: %w", err)
		}
	}()

	if err := yaml.NewEncoder(file).Encode(item); err != nil {
		return fmt.Errorf("marshal item: %w", err)
	} else {
		return nil
	}
}

// SaveToString marshals an item of a registered type to a YAML string.
func (yb *YamlBase) SaveToString(item interface{}) (string, error) {
	builder := &strings.Builder{}

	if err := yaml.NewEncoder(builder).Encode(item); err != nil {
		return "", fmt.Errorf("marshal nexus: %w", err)
	} else {
		return builder.String(), nil
	}
}

//////////////////////////////////////////////////////////////////////////

var typeMatcher = regexp.MustCompile("^" + TypeFieldEscaped + ":\\s+(.+)$")

func getYamlTypeNameAndReset(seeker io.ReadSeeker) (string, error) {
	buffered := bufio.NewReader(seeker)

	for {
		if line, _, err := buffered.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("read line: %w", err)
		} else if matches := typeMatcher.FindStringSubmatch(string(line)); len(matches) < 1 {
			continue
		} else if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return "", fmt.Errorf("seed to beginning of reader: %w", err)
		} else {
			// Trim off any quotes and whitespace.
			return strings.Trim(matches[1], "'\" "), nil
		}
	}

	return "", fmt.Errorf("unable to locate type field")
}
