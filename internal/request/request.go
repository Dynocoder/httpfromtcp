package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	dataString := string(data)

	requestLine, err := parseRequestLine(dataString)
	if err != nil {
		return nil, err
	}

	var request Request

	request.RequestLine = *requestLine

	return &request, nil
}

func parseRequestLine(data string) (*RequestLine, error) {

	lines := strings.Split(data, "\r\n")
	requestLineParts := strings.Split(lines[0], " ")

	var requestLine RequestLine

	// validate request line part count
	if len(requestLineParts) != 3 {
		return nil, errors.New("Request Line should have 3 fields")
	}

	//validate request method
	if !isUppercaseAlphabetic(requestLineParts[0]) {
		return nil, errors.New("Invalid Request Method")
	}
	requestLine.Method = requestLineParts[0]

	requestLine.RequestTarget = requestLineParts[1]

	// validate http request version (only supports 1.1)
	httpVersion := strings.Split(requestLineParts[2], "/")
	if len(httpVersion) != 2 {
		return nil, errors.New("Invalid HTTP Version")
	}

	if httpVersion[1] != "1.1" {
		return nil, errors.New("Invalid HTTP Version")
	}

	requestLine.HttpVersion = httpVersion[1]

	return &requestLine, nil

}

func isUppercaseAlphabetic(data string) bool {
	for _, r := range data {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
