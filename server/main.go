package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	healthy := false
	healthyChan := make(chan bool)

	go func() {
		timeout := time.NewTimer(0)
		for {
			select {
			case <-healthyChan:
				healthy = true
				timeout = time.NewTimer(time.Second * 30)
			case <-timeout.C:
				healthy = false
			}
		}
	}()

	go func() {
		app.Get("/api/healthcheck", func(c *fiber.Ctx) error {
			healthyChan <- true
			return c.SendStatus(200)
		})

		app.Listen(":3000")
	}()

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("%t", healthy)
			}
		}
	}()

	select {
	case <-time.NewTimer(time.Second * 60).C:
		return
	}
}
