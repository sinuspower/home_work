package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
	Notify(string)
}

type Client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	err        io.Writer
	connection net.Conn
	closed     bool
}

var (
	ErrCanNotConnect         = errors.New("can not connect to server")
	ErrCanNotWriteOut        = errors.New("can not write into out")
	ErrCanNotCloseConnection = errors.New("can not close connection")
	ErrCanNotWriteSocket     = errors.New("can not write socket")
)

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address:    address,
		timeout:    timeout,
		in:         in,
		out:        out,
		err:        os.Stderr,
		connection: nil,
		closed:     false,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return ErrCanNotConnect
	}
	c.connection = conn
	c.Notify("...Connected to " + c.address + "\n") // all messages are displayed in Stderr
	return nil
}

func (c *Client) Send() error {
	scanner := bufio.NewScanner(c.in)
	for {
		if !scanner.Scan() {
			c.Notify("...EOF")
			// c.Close() // "basic" test fails!
			break
		}
		if !c.closed {
			_, err := c.connection.Write(append(scanner.Bytes(), '\n'))
			if err != nil {
				return fmt.Errorf("%s: %w", ErrCanNotWriteSocket, err)
			}
		} else {
			break
		}
	}
	return nil
}

func (c *Client) Receive() error {
	scanner := bufio.NewScanner(c.connection)
	for {
		if !scanner.Scan() {
			c.Notify("\n...Connection closed")
			c.closed = true // if connection closed by peer
			break
		}
		_, err := c.out.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return fmt.Errorf("%s: %w", ErrCanNotWriteOut, err)
		}
	}
	return nil
}

func (c *Client) Close() error {
	if c.connection == nil || c.closed {
		return nil
	}
	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("%s: %w", ErrCanNotCloseConnection, err)
	}
	c.closed = true
	return nil
}

func (c *Client) Notify(str string) {
	fmt.Fprintln(c.err, str)
}
