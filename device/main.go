package main

import (
	"encoding/hex"
	"log"
	"os"
	"os/signal"

	"toddler-player/reader"
)

func main() {

	cardReader := reader.NewReader()

	keySignal := make(chan os.Signal, 1)
	signal.Notify(keySignal, os.Interrupt)

	go func() {
		cardReader.Read()
	}()

	go func() {
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
	}()

}
