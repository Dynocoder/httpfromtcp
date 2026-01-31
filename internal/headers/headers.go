package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	read := 0
	doneState := false

	for {
		idx := bytes.Index(data[read:], rn)

		// no rn found, need more data
		if idx == -1 {
			break
		}

		// string end at the start indicates empty line after header (end of headers)
		if idx == 0 {
			doneState = true
			break
		}

		line := []byte(data[read : read+idx])

		name, value, err := parseHeader(line)
		if err != nil {
			return read, false, err
		}

		read += idx + len(rn)
		h[name] = value

	}
	return read, doneState, nil
}

func parseHeader(line []byte) (string, string, error) {
	parts := bytes.SplitN(line, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Bad header format")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("Bad header value format")
	}

	return string(name), string(value), nil
}
