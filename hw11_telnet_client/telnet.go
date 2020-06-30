package main

import (
	"bufio"
	"errors"
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
}

type Client struct {
	Address    string
	Timeout    time.Duration
	In         io.ReadCloser
	Out        io.Writer
	Err        io.Writer
	Connection net.Conn
}

var (
	ErrCanNotConnect         = errors.New("can not connect to server")
	ErrCanNotWriteOut        = errors.New("can not write into out")
	ErrCanNotCloseConnection = errors.New("can not close connection")
	ErrCanNotReadIn          = errors.New("can not read from in")
	ErrCanNotWriteSocket     = errors.New("can not write socket")
	ErrCanNotReadSocket      = errors.New("can not read socket")
)

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		Address:    address,
		Timeout:    timeout,
		In:         in,
		Out:        out,
		Err:        os.Stderr,
		Connection: nil,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.Address, c.Timeout)
	if err != nil {
		return ErrCanNotConnect
	}
	c.Connection = conn
	return nil
}

func (c *Client) Send() error {
	scanner := bufio.NewScanner(c.In)
	if !scanner.Scan() {
		return ErrCanNotReadIn
	}
	bytes := append(scanner.Bytes(), byte('\n'))
	n, err := c.Connection.Write(bytes)
	if err != nil || n == 0 {
		return ErrCanNotWriteSocket
	}
	return nil
}

func (c *Client) Receive() error {
	scanner := bufio.NewScanner(c.Connection)
	if !scanner.Scan() {
		return ErrCanNotReadSocket
	}
	bytes := append(scanner.Bytes(), byte('\n'))
	n, err := c.Out.Write(bytes)
	if err != nil || n == 0 {
		return ErrCanNotWriteOut
	}
	return nil
}

func (c *Client) Close() error {
	if c.Connection == nil {
		return nil
	}
	err := c.Connection.Close()
	if err != nil {
		return ErrCanNotCloseConnection
	}
	return nil
}

// author's solution takes no more than 50 lines
