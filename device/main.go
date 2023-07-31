package main

import (
	"encoding/hex"
	"log"
	"os"
	"os/signal"
)

func main() {

	cardReader := newReader()

	keySignal := make(chan os.Signal, 1)
	signal.Notify(keySignal, os.Interrupt)

	go func() {
		cardReader.read()
	}()

	for {
		select {
		case data := <-cardReader.c:
			log.Println("UID:", hex.EncodeToString(data))
		case <-keySignal:
			log.Println("SIGINT detected, closing toddler-player")
			return
		}
	}
}
