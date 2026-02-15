package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
)

type Server struct {
	open     bool
	listener net.Listener
	handler  handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port uint16, handler handler) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	server := &Server{
		open:     true,
		listener: listener,
		handler:  handler,
	}

	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.open {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {

	defer conn.Close()

	request, err := request.RequestFromReader(conn)
	headers := response.GetDefaultHeaders(0)
	if err != nil {
		HErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		HErr.Write(conn)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	res := s.handler(writer, request)
	if res != nil {
		res.Write(conn)
		return
	}

	headers.Replace("content-length", fmt.Sprintf("%d", len(writer.Bytes())))

	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, headers)
	conn.Write(writer.Bytes())

}

func (he *HandlerError) Write(w io.Writer) error {
	h := response.GetDefaultHeaders(len(he.Message))
	response.WriteStatusLine(w, response.StatusBadRequest)
	response.WriteHeaders(w, h)
	w.Write([]byte(he.Message))
	return nil
}
