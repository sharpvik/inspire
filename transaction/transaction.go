package transaction

import (
	"errors"
	"fmt"
	"net"

	"github.com/bwesterb/go-pow"
	"github.com/sharpvik/inspire/challenge"
	"github.com/sharpvik/inspire/handler"
	"github.com/sharpvik/inspire/message"
)

var Done = message.Message([]byte("DONE"))

type Transaction struct {
	conn      net.Conn
	handle    handler.Handler
	challenge string
}

func New(
	conn net.Conn,
	handle handler.Handler,
	chal challenge.Challenge,
) *Transaction {
	return &Transaction{
		conn:      conn,
		handle:    handle,
		challenge: chal.New(),
	}
}

/*
CLIENT				SERVER
|                   |
|--- REQUEST ------>|
|                   |
|<-- CHALLENGE -----|
|                   |
|--- PROOF -------->|
|                   |
|<-- RESPONSE ------|
|                   |
*/
func (t *Transaction) Handle() {
	request := t.mustRead()
	t.mustSend(message.Message(t.challenge))
	proof := t.mustRead()
	t.mustCheck(proof)
	t.mustSend(t.handle(request))
}

func (t *Transaction) mustCheck(proof []byte) {
	ok, err := pow.Check(t.challenge, string(proof), Done)
	if err != nil {
		panic(fmt.Errorf("check failed: %w", err))
	}
	if !ok {
		panic(errors.New("check failed: PoW"))
	}
}

func (t *Transaction) mustRead() message.Message {
	request, err := message.Read(t.conn)
	if err != nil {
		panic(err)
	}
	return request
}

func (t *Transaction) mustSend(msg message.Message) {
	if err := msg.Send(t.conn); err != nil {
		panic(err)
	}
}
