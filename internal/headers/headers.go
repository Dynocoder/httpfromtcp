package headers

import (
	"bytes"
	"fmt"
	"strings"
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
			read += len(rn)
			break
		}

		line := []byte(data[read : read+idx])

		name, value, err := parseHeader(line)
		if err != nil {
			return read, false, err
		}

		read += idx + len(rn)

		// Header already exists
		// RFC 9110, Section 5.2:
		// https://datatracker.ietf.org/doc/html/rfc9110#name-field-lines-and-combined-fi
		existing := h.Get(name)
		if len(existing) != 0 {
			value = fmt.Sprintf("%s, %s", existing, value)
		}

		h.Set(name, value)

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

	if !isValidToken(name) {
		return "", "", fmt.Errorf("Invalid header token")
	}

	return string(name), string(value), nil
}

func isValidToken(name []byte) bool {
	for _, ch := range name {
		exists := false
		if (ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= 0 && ch <= 9) {
			exists = true
		}

		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			exists = true
		}

		if !exists {
			return false
		}
	}

	return true
}

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

func (h Headers) Set(name string, value string) {
	h[strings.ToLower(name)] = value
}
