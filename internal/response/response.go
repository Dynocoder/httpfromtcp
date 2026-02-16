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

type WriterState int

const (
	WriterStateWaiting WriterState = iota
	WriterStateWrittenStatusLine
	WriterStateWrittenHeaders
	WriterStateWrittenBody
)

type Writer struct {
	io.Writer
	WriterState
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{Writer: writer, WriterState: WriterStateWaiting}
}

func (w *Writer) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.WriterState != WriterStateWaiting {
		return fmt.Errorf("Error: writing status line in wrong state: %v", w.WriterState)
	}

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
	if err != nil {
		return err
	}

	w.WriterState = WriterStateWrittenStatusLine

	return nil
}

func (w *Writer) WriteHeaders(headers *headers.Headers) error {

	if w.WriterState != WriterStateWrittenStatusLine {
		return fmt.Errorf("Error: writing headers in wrong state: %v", w.WriterState)
	}

	h := []byte{}
	for k, v := range *headers {
		h = fmt.Appendf(h, "%s: %s\r\n", k, v)
	}
	h = fmt.Appendf(h, "\r\n")
	_, err := w.Write(h)
	if err != nil {
		return err
	}

	w.WriterState = WriterStateWrittenHeaders

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.WriterState != WriterStateWrittenHeaders {
		return 0, fmt.Errorf("Error: writing headers in wrong state: %v", w.WriterState)
	}

	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, err
	}

	w.WriterState = WriterStateWrittenBody

	return n, nil

}

// Writes to the Writer an HTTP response.
// The method expects WriterState to be `WriterStateWaiting`.
func (w *Writer) ReturnResponse(status StatusCode, h *headers.Headers, p []byte) {
	w.WriteStatusLine(status)
	w.WriteHeaders(h)
	w.WriteBody(p)
}

func GetDefaultHeaders(contentLen int) *headers.Headers {

	h := headers.NewHeaders()

	h.Set("content-length", strconv.Itoa(contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")

	return &h
}
