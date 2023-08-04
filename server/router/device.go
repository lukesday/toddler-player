package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) UseDevice() {
	r.Conn.SetStatus("health", "offline")
	healthyChan := make(chan bool)

	go func() {
		timeout := time.NewTimer(0)
		for {
			select {
			case status := <-healthyChan:
				if status {
					r.Conn.SetStatus("health", "online")
					timeout = time.NewTimer(time.Second * 30)
				} else {
					r.Conn.SetStatus("health", "offline")
					timeout = time.NewTimer(0)
				}
			case <-timeout.C:
				r.Conn.SetStatus("health", "offline")
			}
		}
	}()

	r.App.Put("/api/device/sethealth", func(c *fiber.Ctx) error {
		healthyChan <- true
		return c.SendStatus(200)
	})

	r.App.Get("/api/device/healthcheck", func(c *fiber.Ctx) error {
		healthStatus, err := r.Conn.CheckStatus("health")
		if err == nil && healthStatus == "online" {
			return c.JSON(fiber.Map{"status": "online"})
		}
		return c.JSON(fiber.Map{"status": "offline"})
	})
}
