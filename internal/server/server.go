package server

import (
	"fmt"
	"github.com/peeta98/httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{listener: listener}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		err := s.listener.Close()
		return err
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
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

	// Write the status line first ("HTTP/1.1 200 OK")
	response.WriteStatusLine(conn, response.StatusCodeOK)
	// After that, write the headers
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))

	return
}
