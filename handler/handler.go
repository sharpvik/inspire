package handler

import (
	"github.com/sharpvik/purr/message"
)

// Any function that accepts a Message and returns a Message is a Handler.
type Handler func(message.Message) message.Message

// Echo is a simple Handler that responds with the same Message as it received.
func Echo(msg message.Message) message.Message { return msg }

// EchoIfNil is a utility function used to fill the Server.handler field with
// a sensible default (Echo handler) during construction.
func EchoIfNil(h Handler) Handler {
	if h != nil {
		return h
	}
	return Echo
}
