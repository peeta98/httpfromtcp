package server

import (
	"bytes"
	"fmt"
	"github.com/peeta98/httpfromtcp/internal/request"
	"github.com/peeta98/httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		handler:  handler,
		listener: listener,
	}

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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerError := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		handlerError.Write(conn)
		return
	}

	buf := new(bytes.Buffer)

	handlerError := s.handler(buf, req)
	if handlerError != nil {
		handlerError.Write(conn)
		return
	}

	// Write the status line first ("HTTP/1.1 200 OK")
	response.WriteStatusLine(conn, response.StatusCodeOK)

	defaultHeaders := response.GetDefaultHeaders(buf.Len())
	// After that, write the headers
	response.WriteHeaders(conn, defaultHeaders)
	conn.Write(buf.Bytes())
}
