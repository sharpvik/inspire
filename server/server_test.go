package server

import (
	"testing"
	"time"

	"github.com/sharpvik/purr/client"
	"github.com/sharpvik/purr/message"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	server := New(nil)
	msg := message.Message([]byte("hello"))
	response := server.Handle(msg)
	assert.Equal(t, msg, response)
}

func TestClient(t *testing.T) {
	const addr = "localhost:8910"
	server := New(nil)
	msg := message.Message([]byte("hello"))
	peer := client.New(addr)
	go server.ListenAndServe(addr)
	time.Sleep(time.Second) // allow server to spin up
	response, err := peer.Send(msg)
	assert.NoError(t, err)
	assert.Equal(t, msg, response)
}
