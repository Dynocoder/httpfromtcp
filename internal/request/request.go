package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"
)

type ParserState int

const (
	initialized ParserState = iota
	done
)
const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState == initialized {
		line, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		} else if n == 0 {
			return 0, nil // did not receive CRLF, keep going
		} else if n > 0 {
			r.RequestLine = *line
			r.ParserState = done
			return n, nil
		}

	} else if r.ParserState == done {
		return 0, errors.New("error: trying to read data in done state")
	} else {
		return 0, errors.New("error: Unknown state")
	}

	return 0, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)

	request := &Request{
		ParserState: initialized,
	}

	readToIndex := 0

	for request.ParserState != done {

		// if buffer is full, grow it
		if readToIndex == cap(buf) {
			newBuf := make([]byte, readToIndex, cap(buf)*2)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:cap(buf)])
		if n > 0 {
			readToIndex += n
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		consumed, perr := request.parse(buf)
		if perr != nil {
			return nil, perr
		}

		if consumed > 0 {
			copy(buf, buf[consumed:readToIndex])
			readToIndex -= consumed
		}
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {

	idx := bytes.Index(data, []byte("\r\n"))

	if idx == -1 {
		return &RequestLine{}, 0, nil
	}

	line := string(data[:idx])
	n := len(line)
	requestLineParts := strings.Split(line, " ")

	var requestLine RequestLine

	// validate request line part count
	if len(requestLineParts) != 3 {
		return nil, n, errors.New("Request Line should have 3 fields")
	}

	//validate request method
	if !isUppercaseAlphabetic(requestLineParts[0]) {
		return nil, n, errors.New("Invalid Request Method")
	}
	requestLine.Method = requestLineParts[0]

	requestLine.RequestTarget = requestLineParts[1]

	// validate http request version (only supports 1.1)
	httpVersion := strings.Split(requestLineParts[2], "/")
	if len(httpVersion) != 2 {
		return nil, n, errors.New("Invalid HTTP Version")
	}

	if httpVersion[1] != "1.1" {
		return nil, n, errors.New("Invalid HTTP Version")
	}

	requestLine.HttpVersion = httpVersion[1]

	return &requestLine, n, nil

}

func isUppercaseAlphabetic(data string) bool {
	for _, r := range data {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
