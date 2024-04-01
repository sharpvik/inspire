package message

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Message []byte

func Read(rd io.Reader) (Message, error) {
	r := bufio.NewReader(rd)

	var length int32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	msg := make([]byte, length)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg Message) Send(w io.Writer) error {
	_, err := w.Write(msg)
	return err
}
