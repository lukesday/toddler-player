package main

import (
	"time"

	"toddler-player/server/database"
	"toddler-player/server/router"
)

func main() {
	conn := database.InitialiseDatabase()
	router.InitialiseRouter(&conn)

	select {
	case <-time.NewTimer(time.Minute * 60).C:
		return
	}
}
