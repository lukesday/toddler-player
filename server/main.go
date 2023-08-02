package main

import (
	"os"
	"time"

	"toddler-player/server/database"
	"toddler-player/server/router"

	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	if env == "dev" {
		godotenv.Load()
	}

	conn := database.InitialiseDatabase()
	router.InitialiseRouter(&conn)

	select {
	case <-time.NewTimer(time.Minute * 60).C:
		return
	}
}
