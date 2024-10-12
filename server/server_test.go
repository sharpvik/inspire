package server_test

import (
	"testing"

	"github.com/sharpvik/inspire/challenge"
	"github.com/sharpvik/inspire/message"
	"github.com/sharpvik/inspire/server"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	server := server.New(nil, challenge.WithDifficulty(32))
	msg := message.Message([]byte("hello"))
	response := server.Handle(msg)
	assert.Equal(t, msg, response)
}
