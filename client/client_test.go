package client_test

import (
	"testing"
	"time"

	"github.com/sharpvik/inspire/challenge"
	"github.com/sharpvik/inspire/client"
	"github.com/sharpvik/inspire/message"
	"github.com/sharpvik/inspire/server"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	const addr = "localhost:8910"
	server := server.New(nil, challenge.WithDifficulty(32))
	msg := message.Message([]byte("hello"))
	peer := client.New(addr)
	go server.ListenAndServe(addr)
	time.Sleep(time.Second) // allow server to spin up
	response, err := peer.Send(msg)
	require.NoError(t, err)
	require.Equal(t, msg, response)
}
