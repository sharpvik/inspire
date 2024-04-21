package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/sharpvik/purr/handler"
	"github.com/sharpvik/purr/message"
)

//	A Server's life ain't fun as
//	Most Clients just want something done.
//	They dial and Send, while I must pretend...
//	Submerged into reading a Message,
//	I wish it had never been sent!
//
// Server is just sitting there on the network waiting for Clients to connect.
// Each Message is handled by the Handler given to Server during construction.
type Server struct {
	/* Configurable */

	handler handler.Handler

	/* Internal */

	listener net.Listener
}

// Construct a new Server instance with a given handler.
func New(h handler.Handler) *Server {
	return &Server{
		handler: handler.EchoIfNil(h), // sensible default
	}
}

// Attach to the address, listen for client connections, and serve them.
func (s *Server) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	return s.acceptConnections()
}

// Handling an incoming Message is simple - just call the Handler.
func (s *Server) Handle(msg message.Message) message.Message {
	return s.handler(msg)
}

// This loop is run forever until the server is stopped.
func (s *Server) acceptConnections() error {
	for {
		if err := s.acceptAndHandleConnection(); err != nil {
			return err
		}
	}
}

// Take the next connection request from the listener, accept it (thus
// establishing a new connection) and handle in a separate thread of execution.
func (s *Server) acceptAndHandleConnection() error {
	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}
	go s.handleConnection(conn)
	return nil
}

// To handle a connection we must simply handle incoming Messages one-by-one
// until the connection is closed.
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

// Handle request by reading a Message, coming up with a response, and sending
// it back to the client.
func (s *Server) handleRequest(conn net.Conn) error {
	request, err := message.Read(conn)
	if err != nil {
		return err
	}
	return s.Handle(request).Send(conn)
}
