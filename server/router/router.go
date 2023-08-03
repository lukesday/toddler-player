package router

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

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

		app.All("/api/*", func(c *fiber.Ctx) error {
			conn.LogRequest(c.Request().URI().String(), string(c.Request().Header.Method()))
			return c.Next()
		})

		app.Get("/api/logs", func(c *fiber.Ctx) error {
			return c.JSON(conn.GetLogs())
		})

		app.Put("/api/device/sethealth", func(c *fiber.Ctx) error {
			healthyChan <- true
			return c.SendStatus(200)
		})

		app.Get("/api/device/healthcheck", func(c *fiber.Ctx) error {
			healthStatus, err := conn.CheckStatus("health")
			if err == nil && healthStatus == "online" {
				return c.JSON(fiber.Map{"status": "online"})
			}
			return c.JSON(fiber.Map{"status": "offline"})
		})

		app.Post("/api/nfc/:uid", func(c *fiber.Ctx) error {
			nfcTag := database.NfcTag{}
			if tagErr := conn.GetTag(c.Params("uid"), nfcTag); tagErr == nil {
				automation := database.Automation{}

				if err := conn.GetAutomation(nfcTag, automation); err != nil {
					log.Println(err)
				} else {
					// Trigger automation
					log.Println(automation)
				}
			} else if tagErr == gorm.ErrRecordNotFound {
				conn.Enroll(c.Params("uid"))
			}

			return c.SendStatus(200)
		})

		app.Post("/api/automation", func(c *fiber.Ctx) error {
			payload := struct {
				NfcTag   string `json:"nfcTag"`
				DeviceId string `json:"deviceId"`
				MediaId  string `json:"mediaId"`
			}{}

			if err := c.BodyParser(&payload); err != nil {
				log.Println("error = ", err)
				return err
			}

			nfcTag := database.NfcTag{}

			if err := conn.GetTag(payload.NfcTag, nfcTag); err != nil {
				log.Println("error = ", err)
				return err
			}

			automation := database.Automation{}

			if err := conn.GetAutomation(nfcTag, automation); err == nil {
				return c.SendStatus(409)
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				conn.CreateAutomation(nfcTag, payload.DeviceId, payload.MediaId)
				return c.SendStatus(200)
			}

			return c.SendStatus(400)
		})

		app.Listen(":3000")
	}()
}
