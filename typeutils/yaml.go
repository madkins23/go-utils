package typeutils

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

//////////////////////////////////////////////////////////////////////////

var typeMatcher = regexp.MustCompile("^" + TypeFieldEscaped + ":\\s+(.+)$")

func GetYamlTypeNameAndReset(seeker io.ReadSeeker) (string, error) {
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
