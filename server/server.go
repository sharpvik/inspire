package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/sharpvik/purr/handler"
	"github.com/sharpvik/purr/message"
)

type Server struct {
	/* Configurable */

	handler handler.Handler

	/* Internal */

	listener net.Listener
}

func New(h handler.Handler) *Server {
	return &Server{
		handler: handler.EchoIfNil(h), // sensible default
	}
}

func (s *Server) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	return s.acceptConnections()
}

func (s *Server) Handle(msg message.Message) message.Message {
	return s.handler(msg)
}

func (s *Server) acceptConnections() error {
	for {
		if err := s.acceptAndHandleConnection(); err != nil {
			return err
		}
	}
}

func (s *Server) acceptAndHandleConnection() error {
	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}
	go s.handleConnection(conn)
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.recoverFromConnectionPanic(conn)
	for {
		if err := s.handleRequest(conn); err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) recoverFromConnectionPanic(conn net.Conn) {
	if something := recover(); something != nil {
		if err, ok := something.(error); ok {
			log.Println(errors.Join(connectionHandlingError(conn), err))
			conn.Close()
		}
	}
}

func connectionHandlingError(conn net.Conn) error {
	return fmt.Errorf(
		"encountered a panic during connection handling (addr %s | %s)",
		conn.LocalAddr(), conn.RemoteAddr(),
	)
}

func (s *Server) handleRequest(conn net.Conn) error {
	request, err := message.Read(conn)
	if err != nil {
		return err
	}
	return s.Handle(request).Send(conn)
}
