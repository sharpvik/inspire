package client

import (
	"net"

	"github.com/sharpvik/purr/message"
)

// Client is happy to connect to the Server and exchange some data with it.
type Client struct {
	/* Configurable */

	serverAddress string

	/* Internal */

	conn net.Conn
}

// Constructs a fresh-new Client with the specified Server address.
func New(serverAddress string) *Client {
	return &Client{
		serverAddress: serverAddress,
	}
}

// It's not about the money. It's about sending the Message. The Client makes it
// as simple as it could be. Just give it the Message you'd like to send and it
// will connect to the Server (if necessary) and pass it over like a good boy.
func (c *Client) Send(msg message.Message) (message.Message, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	if err := msg.Send(c.conn); err != nil {
		return nil, err
	}
	return message.Read(c.conn)
}

// We don't have to redial the Server if we already have a connection. This
// method ensures that we only dial when necessary.
func (c *Client) connect() error {
	if c.isAlreadyConnected() {
		return nil
	}
	return c.dial()
}

// Dialing is also known as establishing a connection.
func (c *Client) dial() error {
	conn, err := net.Dial("tcp", c.serverAddress)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) isAlreadyConnected() bool {
	return c.conn != nil
}
