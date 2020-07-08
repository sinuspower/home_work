package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		return
	}

	if err := generate(args[1]); err != nil {
		log.Fatal(err)
	}
}
