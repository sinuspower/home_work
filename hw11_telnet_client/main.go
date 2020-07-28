package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "server connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		return
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	client.Notify("...Connected to " + address + "\n")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	sigDone := make(chan struct{})
	go func() {
		<-c
		signal.Stop(c)
		sigDone <- struct{}{}
	}()

	readDone := make(chan struct{})
	go func() {
		read(client)
		readDone <- struct{}{}
	}()

	writeDone := make(chan struct{})
	go func() {
		write(client)
		writeDone <- struct{}{}
	}()

	select {
	case <-sigDone:
		client.Notify("\n...Cancelled")
	case <-readDone:
		client.Notify("\n...Connection closed by peer")
	case <-writeDone:
		client.Notify("\n...EOF")
	}
}

func read(client TelnetClient) {
	if err := client.Receive(); err != nil {
		log.Fatalf("%s: %s", "Error receiving: ", err)
	}
}

func write(client TelnetClient) {
	if err := client.Send(); err != nil {
		log.Fatalf("%s: %s", "Error sending: ", err)
	}
}
