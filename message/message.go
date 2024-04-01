package message

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Message []byte

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

func (msg Message) Send(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint32(len(msg))); err != nil {
		return err
	}
	_, err := w.Write(msg)
	return err
}
