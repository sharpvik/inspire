package handler

import (
	"github.com/sharpvik/purr/message"
)

type Handler func(message.Message) message.Message

func Echo(msg message.Message) message.Message {
	return msg
}

func EchoIfNil(h Handler) Handler {
	if h != nil {
		return h
	}
	return Echo
}
