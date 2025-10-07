package main

import (
	"log"
	"time"

	"toddler-player/server/database"
	"toddler-player/server/router"

	"github.com/joho/godotenv"
)

func main() {
	//env := os.Getenv("ENV")
	//if env == "dev" {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	//}

	conn := database.InitialiseDatabase()
	router.InitialiseRouter(&conn)

	select {
	case <-time.NewTimer(time.Minute * 60).C:
		return
	}
}
