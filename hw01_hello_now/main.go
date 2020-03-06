package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

const (
	layout           string = "2006-01-02 15:04:05 +0000 UTC" // time format
	host             string = "0.beevik-ntp.pool.ntp.org"     // host for getting exact time
	exitCodeNTPError int    = 1
)

func main() {
	now := time.Now()
	fmt.Printf("current time: %s\n", now.Format(layout))
	if exact, err := ntp.Time(host); err != nil {
		fmt.Fprintf(os.Stderr, "error getting exact time: %v\n\t", err)
		os.Exit(exitCodeNTPError)
	} else {
		fmt.Printf("exact time: %s\n", exact.Format(layout))
	}
}
