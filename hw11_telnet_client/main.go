package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
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

	var wg sync.WaitGroup
	wg.Add(2)

	go read(client, &wg)
	go write(client, &wg)

	wg.Wait()
}

func read(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := client.Receive(); err != nil {
		log.Fatal(err)
	}
}

func write(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := client.Send(); err != nil {
		log.Fatal(err)
	}
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
