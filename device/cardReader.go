package main

import (
	"log"
	"time"

	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

type reader struct {
	c chan []byte
}

func newReader() reader {
	thisReader := reader{
		c: make(chan []byte),
	}
	return thisReader
}

func (r *reader) read() {
	log.Printf("Loading read")

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

	loopTimer := time.NewTimer(0)

	// Stopping timer, flagging reader thread as timed out
	defer func() {
		close(r.c)
	}()

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

			r.c <- uid
		}
	}
}
