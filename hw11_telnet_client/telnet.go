package main

import (
	"bufio"
	"context"
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
	Cancel()
	Notify(string)
}

type Client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	err        io.Writer
	connection net.Conn
	context    context.Context
	cancel     func()
}

var (
	ErrCanNotConnect         = errors.New("can not connect to server")
	ErrCanNotWriteOut        = errors.New("can not write into out")
	ErrCanNotCloseConnection = errors.New("can not close connection")
	ErrCanNotWriteSocket     = errors.New("can not write socket")
)

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	client := Client{
		address:    address,
		timeout:    timeout,
		in:         in,
		out:        out,
		err:        os.Stderr,
		connection: nil,
		context:    context.Background(),
	}
	ctx, cancel := context.WithCancel(client.context)
	client.context = ctx
	client.cancel = cancel
	return &client
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return ErrCanNotConnect
	}
	c.connection = conn
	c.Notify("...Connected to " + c.address) // all messages are displayed in Stderr
	return nil
}

func (c *Client) Send() error {
	scanner := bufio.NewScanner(c.in)
SEND:
	for {
		select {
		case <-c.context.Done():
			c.Notify("...Connection was closed by peer")
			break SEND
		default:
			if !scanner.Scan() {
				c.Notify("...EOF")
				break SEND
			}
			str := scanner.Text()
			_, err := c.connection.Write([]byte(fmt.Sprintf("%s\n", str)))
			if err != nil {
				return ErrCanNotWriteSocket
			}
		}
	}
	return nil
}

func (c *Client) Receive() error {
	scanner := bufio.NewScanner(c.connection)
RECEIVE:
	for {
		select {
		case <-c.context.Done():
			break RECEIVE
		default:
			if !scanner.Scan() {
				c.Cancel() // connection closed, stop sending
				break RECEIVE
			}
			str := scanner.Text()
			_, err := c.out.Write([]byte(fmt.Sprintf("%s\n", str)))
			if err != nil {
				return ErrCanNotWriteOut
			}
		}
	}
	return nil
}

func (c *Client) Close() error {
	if c.connection == nil {
		return nil
	}
	err := c.connection.Close()
	if err != nil {
		return ErrCanNotCloseConnection
	}
	return nil
}

func (c *Client) Cancel() {
	c.cancel()
}

func (c *Client) Notify(str string) {
	fmt.Fprintln(c.err, str)
}
