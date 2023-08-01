package main

import (
	"time"

	"toddler-player/server/router"
)

func main() {
	router.InitialiseRouter()

	select {
	case <-time.NewTimer(time.Second * 60).C:
		return
	}
}
