package main

import (
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"time"

	"net/http"
	"toddler-player/device/reader"
)

func main() {

	const (
		API_BASE_URL string = "http://192.168.1.109:3000"
	)

	client := &http.Client{}

	cardReader := reader.NewReader()

	keySignal := make(chan os.Signal, 1)
	signal.Notify(keySignal, os.Interrupt)

	// Initialise card reader
	go cardReader.Read()

	// Health check for API
	go func() {
		loopTimer := time.NewTimer(0)

		for {
			<-loopTimer.C
			// Send healthy HTTP request to endpoint instead
			log.Println("healthy")
			req, _ := http.NewRequest(http.MethodPut, API_BASE_URL+"/api/device/sethealth", nil)
			resp, err := client.Do(req)
			if err != nil {
				log.Println("Error sending healthy HTTP request:", err)
			} else {
				resp.Body.Close()
				if resp.StatusCode != 200 {
					log.Println("Error sending healthy HTTP request:", resp.StatusCode)
				}
			}
			loopTimer = time.NewTimer(time.Second * 10)
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
