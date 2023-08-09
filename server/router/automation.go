package router

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"toddler-player/server/database"
)

type AutomationPayload struct {
	NfcTagUid string `json:"nfcTagUid"`
	DeviceId  string `json:"deviceId"`
	MediaId   string `json:"mediaId"`
}

func (r *Router) UseAutomation() {
	r.App.Get("/api/automation", func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return c.SendStatus(400)
		}
		user := localUser.(database.User)

		automations := []database.Automation{}
		if err := r.Conn.ListAutomations(user.UserID, &automations); err != nil {
			return err
		}
		return c.JSON(automations)
	})

	r.App.Get("/api/automation/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.SendStatus(400)
		}

		automation := database.Automation{}
		if err := r.Conn.GetAutomation(uint(id), &automation); err != nil {
			if errors.Is(err, fiber.ErrNotFound) {
				return c.SendStatus(404)
			}
			return c.SendStatus(500)
		}

		return c.JSON(automation)
	})

	r.App.Post("/api/automation", func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return c.SendStatus(400)
		}
		user := localUser.(database.User)

		payload := AutomationPayload{}

		if err := c.BodyParser(&payload); err != nil {
			log.Println("error = ", err)
			return err
		}

		nfcTag := database.NfcTag{}

		if err := r.Conn.GetTag(payload.NfcTagUid, &nfcTag); err != nil {
			log.Println("error = ", err)
			return err
		}

		automation := database.Automation{}

		if err := r.Conn.GetAutomationByNfcTag(nfcTag, &automation); err == nil {
			return c.SendStatus(409)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Conn.CreateAutomation(nfcTag, payload.DeviceId, payload.MediaId, user)
			return c.SendStatus(200)
		}

		return c.SendStatus(400)
	})
}
