package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
)

type Server struct {
	open     bool
	listener net.Listener
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port uint16, handler Handler) (*Server, error) {

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
	return s.listener.Close()
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

	headers := response.GetDefaultHeaders(0)
	writer := response.NewWriter(conn)
	request, err := request.RequestFromReader(conn)
	if err != nil {
		writer.ReturnResponse(response.StatusBadRequest, headers, []byte{})
		return
	}

	s.handler(writer, request)

}
