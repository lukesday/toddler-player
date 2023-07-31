// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"time"

	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

func main() {
	log.Printf("test")

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Using SPI as an example. See package "periph.io/x/conn/v3/spi/spireg" for more details.
	p, err := spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	rfid, err := mfrc522.NewSPI(p, rpi.P1_13, rpi.P1_11)
	if err != nil {
		log.Fatal(err)
	}

	// Idling device on exit.
	defer rfid.Halt()

	// Setting the antenna signal strength.
	rfid.SetAntennaGain(5)

	timedOut := false
	cb := make(chan []byte)
	timer := time.NewTimer(10 * time.Minute)

	keySignal := make(chan os.Signal, 1)
	signal.Notify(keySignal, os.Interrupt)

	// Stopping timer, flagging reader thread as timed out
	defer func() {
		timer.Stop()
		timedOut = true
		close(cb)
	}()

	go func() {
		log.Printf("Started %s", rfid.String())

		loopTimer := time.NewTimer(0)

		for {
			// Trying to read card UID.

			select {
			case <-loopTimer.C:
				uid, err := rfid.ReadUID(10 * time.Second)

				// If main thread timed out just exiting.
				if timedOut {
					return
				}

				// Some devices tend to send wrong data while RFID chip is already detected
				// but still "too far" from a receiver.
				// Especially some cheap CN clones which you can find on GearBest, AliExpress, etc.
				// This will suppress such errors.
				if err != nil {
					continue
				}

				loopTimer = time.NewTimer(time.Second * 1)

				cb <- uid
			}
		}
	}()

	for {
		select {
		case <-timer.C:
			log.Fatal("Didn't receive device data")
			return
		case data := <-cb:
			log.Println("UID:", hex.EncodeToString(data))
		case <-keySignal:
			log.Println("SIGINT detected, closing toddler-player")
			return
		}
	}
}
