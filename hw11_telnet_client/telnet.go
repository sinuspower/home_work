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
}

var (
	ErrCanNotConnect         = errors.New("can not connect to server")
	ErrCanNotWriteOut        = errors.New("can not write into out")
	ErrCanNotCloseConnection = errors.New("can not close connection")
)

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address:    address,
		timeout:    timeout,
		in:         in,
		out:        out,
		err:        os.Stderr,
		connection: nil,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return ErrCanNotConnect
	}
	c.connection = conn
	return nil
}

func (c *Client) Send() error {
	_, err := io.Copy(c.connection, c.in)
	return err
}

func (c *Client) Receive() error {
	scanner := bufio.NewScanner(c.connection)
	for {
		if !scanner.Scan() {
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
	if c.connection == nil {
		return nil
	}
	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("%s: %w", ErrCanNotCloseConnection, err)
	}
	return nil
}

// Notify writes a message into errors stream (Stderr).
func (c *Client) Notify(str string) {
	fmt.Fprintln(c.err, str)
}
