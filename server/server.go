package server

import (
	"io"
	"log"
	"net"

	"github.com/sharpvik/inspire/challenge"
	"github.com/sharpvik/inspire/handler"
	"github.com/sharpvik/inspire/message"
	"github.com/sharpvik/inspire/transaction"
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

	handle    handler.Handler
	challenge challenge.Challenge

	/* Internal */

	listener net.Listener
}

// Construct a new Server instance with a given handler.
func New(h handler.Handler, c challenge.Challenge) *Server {
	return &Server{
		handle:    handler.EchoIfNil(h), // sensible default
		challenge: c,
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
	return s.handle(msg)
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

func (s *Server) handleConnection(conn net.Conn) {
	defer s.recoverFromConnectionPanic(conn)
	for {
		transaction.New(conn, s.handle, s.challenge).Handle()
	}
}

func (s *Server) recoverFromConnectionPanic(conn net.Conn) {
	conn.Close()

	something := recover()
	if something == nil {
		return
	}

	err, ok := something.(error)
	if !ok || err == io.EOF {
		return
	}

	log.Println(err)
}
