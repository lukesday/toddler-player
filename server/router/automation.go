package router

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"toddler-player/server/database"
)

type AutomationPayload struct {
	NfcTag   string `json:"nfcTag"`
	DeviceId string `json:"deviceId"`
	MediaId  string `json:"mediaId"`
}

func (r *Router) UseAutomation() {
	r.App.Post("/api/automation", func(c *fiber.Ctx) error {
		payload := AutomationPayload{}

		if err := c.BodyParser(&payload); err != nil {
			log.Println("error = ", err)
			return err
		}

		nfcTag := database.NfcTag{}

		if err := r.Conn.GetTag(payload.NfcTag, nfcTag); err != nil {
			log.Println("error = ", err)
			return err
		}

		automation := database.Automation{}

		if err := r.Conn.GetAutomation(nfcTag, automation); err == nil {
			return c.SendStatus(409)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Conn.CreateAutomation(nfcTag, payload.DeviceId, payload.MediaId)
			return c.SendStatus(200)
		}

		return c.SendStatus(400)
	})
}
