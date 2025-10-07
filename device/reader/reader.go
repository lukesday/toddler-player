package reader

import (
	"log"
	"time"

	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
	"periph.io/x/host/v3"
)

type Reader struct {
	C chan []byte
}

func NewReader() Reader {
	thisReader := Reader{
		C: make(chan []byte),
	}
	return thisReader
}

func (r *Reader) Read() {
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

	// Access GPIO pins by their names or numbers
	resetPin := gpioreg.ByName("25") // GPIO25 (was rpi.P1_22)
	if resetPin == nil {
		log.Fatal("Failed to find GPIO25")
	}

	// irqPin := gpioreg.ByName("24") // GPIO24 (was rpi.P1_18)
	// if irqPin == nil {
	// 	log.Fatal("Failed to find GPIO24")
	// }

	rfid, err := mfrc522.NewSPI(p, resetPin, nil)
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
		close(r.C)
	}()

	log.Printf("Started %s", rfid.String())

	for {
		// Trying to read card UID.
		<-loopTimer.C
		uid, err := rfid.ReadUID(10 * time.Second)

		loopTimer = time.NewTimer(time.Second * 1)

		// Some devices tend to send wrong data while RFID chip is already detected
		// but still "too far" from a receiver.
		// Especially some cheap CN clones which you can find on GearBest, AliExpress, etc.
		// This will suppress such errors.
		if err != nil {
			continue
		}

		r.C <- uid
	}
}
