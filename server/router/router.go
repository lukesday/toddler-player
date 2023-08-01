package router

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"toddler-player/server/database"
)

func InitialiseRouter(conn *database.DatabaseConnection) {
	app := fiber.New()

	conn.SetStatus("health", "offline")
	healthyChan := make(chan bool)

	go func() {
		timeout := time.NewTimer(0)
		for {
			select {
			case status := <-healthyChan:
				if status {
					conn.SetStatus("health", "online")
					timeout = time.NewTimer(time.Second * 30)
				} else {
					conn.SetStatus("health", "offline")
					timeout = time.NewTimer(0)
				}
			case <-timeout.C:
				conn.SetStatus("health", "offline")
			}
		}
	}()

	go func() {
		app.Put("/api/sethealth", func(c *fiber.Ctx) error {
			healthyChan <- true
			return c.SendStatus(200)
		})

		app.Get("/api/healthcheck", func(c *fiber.Ctx) error {
			healthStatus, err := conn.CheckStatus("health")
			if err == nil && healthStatus == "online" {
				return c.JSON(fiber.Map{"status": "online"})
			}
			return c.JSON(fiber.Map{"status": "offline"})
		})

		app.Post("/api/enroll/:uid", func(c *fiber.Ctx) error {
			if conn.Exists(c.Params("uid")) {
				return c.SendStatus(409)
			}

			conn.Enroll(c.Params("uid"))
			return c.SendStatus(200)
		})

		app.Listen(":3000")
	}()
}
