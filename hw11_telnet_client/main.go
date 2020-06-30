package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	// do not rush to throw context down, think if it is useful with blocking operation?
	var timeoutString string
	flag.StringVar(&timeoutString, "timeout", "10s", "server connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		return
	}

	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		log.Fatal(err)
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("...Connected to %s\n", address)
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("...Connection closed")
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	// write
	go func() {
		var err error
		for err == nil {
			err = client.Send()
		}
		defer wg.Done()
	}()

	// read
	go func() {
		var err error
		for err == nil {
			err = client.Receive()
		}
		defer wg.Done()
	}()

	wg.Wait()
}
