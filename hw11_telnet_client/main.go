package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	go func() {
		for range c {
			client.Close() // CTRL+C
			break
		}
	}()

	go func() {
		read(client)
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		write(client)
	}()

	wg.Wait()
}

func read(client TelnetClient) {
	if err := client.Receive(); err != nil {
		log.Fatal(err)
	}
}

func write(client TelnetClient) {
	if err := client.Send(); err != nil {
		log.Fatal(err)
	}
}
