package main

import (
	"log"
	"os"
	"time"
)

func main() {
	host, _ := os.Hostname()
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			log.Println(host)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
