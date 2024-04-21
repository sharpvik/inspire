package message

import (
	"bufio"
	"encoding/binary"
	"io"
)

// Message is just a bunch of bytes. PURR does not make any assumptions as to
// the kinds of data its users might be exchanging. We're flexible like that.
type Message []byte

// Read a Message from a Reader (which is practically an incoming byte stream).
func Read(rd io.Reader) (Message, error) {
	r := bufio.NewReader(rd)

	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	msg := make([]byte, length)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// Send a Message to a Writer (which is practically an outgoing byte stream).
func (msg Message) Send(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint32(len(msg))); err != nil {
		return err
	}
	_, err := w.Write(msg)
	return err
}
