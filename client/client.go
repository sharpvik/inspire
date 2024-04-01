package client

import (
	"net"

	"github.com/sharpvik/purr/message"
)

type Client struct {
	/* Configurable */

	addr string

	/* Internal */

	conn net.Conn
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Send(msg message.Message) (message.Message, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	if err := msg.Send(c.conn); err != nil {
		return nil, err
	}
	return message.Read(c.conn)
}

func (c *Client) connect() error {
	if c.alreadyConnected() {
		return nil
	}
	return c.dial()
}

func (c *Client) dial() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) alreadyConnected() bool {
	return c.conn != nil
}
