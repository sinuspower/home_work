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
	args = args[1:]

	envDir := args[0]
	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	cmd := args[1:]
	if returnCode := RunCmd(cmd, env); returnCode != 0 {
		os.Exit(returnCode)
	}
}
