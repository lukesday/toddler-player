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

func Read() {
	log.Printf("Loading Read")

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

	cb := make(chan []byte)
	loopTimer := time.NewTimer(0)

	keySignal := make(chan os.Signal, 1)
	signal.Notify(keySignal, os.Interrupt)

	// Stopping timer, flagging reader thread as timed out
	defer func() {
		close(cb)
	}()

	go func() {
		log.Printf("Started %s", rfid.String())

		for {
			// Trying to read card UID.

			select {
			case <-loopTimer.C:
				uid, err := rfid.ReadUID(10 * time.Second)

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
		case data := <-cb:
			log.Println("UID:", hex.EncodeToString(data))
		case <-keySignal:
			log.Println("SIGINT detected, closing toddler-player")
			return
		}
	}
}
