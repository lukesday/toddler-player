package main

import (
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"time"

	"toddler-player/reader"
)

func main() {

	cardReader := reader.NewReader()

	keySignal := make(chan os.Signal, 1)
	signal.Notify(keySignal, os.Interrupt)

	// Initialise card reader
	go cardReader.Read()

	// Health check for API
	go func() {
		loopTimer := time.NewTimer(0)

		for {
			select {
			case <-loopTimer.C:
				// Send healthy HTTP request to endpoint instead
				log.Println("healthy")
				loopTimer = time.NewTimer(time.Second * 10)
			}
		}
	}()

	for {
		select {
		case data := <-cardReader.C:
			log.Println("UID:", hex.EncodeToString(data))
			// Send HTTP request to server with UID
		case <-keySignal:
			log.Println("SIGINT detected, closing toddler-player")
			os.Exit(1)
		}
	}

}
