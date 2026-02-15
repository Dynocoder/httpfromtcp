package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	out := []byte("")
	switch statusCode {
	case StatusOK:
		out = []byte("HTTP/1.1 200 OK \r\n")
	case StatusBadRequest:
		out = []byte("HTTP/1.1 400 Bad Request \r\n")
	case StatusInternalServerError:
		out = []byte("HTTP/1.1 500 Internal Server Error \r\n")
	default:
		return fmt.Errorf("Error: Unrecognized status code")
	}

	_, err := w.Write(out)
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {

	h := headers.NewHeaders()

	h.Set("content-length", strconv.Itoa(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")

	return &h
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {

	h := []byte{}
	for k, v := range *headers {
		h = fmt.Appendf(h, "%s: %s\r\n", k, v)
	}
	h = fmt.Appendf(h, "\r\n")
	_, err := w.Write(h)
	return err
}
