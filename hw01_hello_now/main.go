package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const (
	layout string = "2006-01-02 15:04:05 +0000 UTC" // time format
	host   string = "0.beevik-ntp.pool.ntp.org"     // host for getting exact time
)

func main() {
	now := time.Now()
	exact, err := ntp.Time(host)
	if err != nil {
		log.Fatalf("error getting exact time: %v", err)
	}
	fmt.Printf("current time: %s\nexact time: %s\n", now.Format(layout), exact.Format(layout))
}
