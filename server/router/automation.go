package router

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"toddler-player/server/database"
	"toddler-player/server/spotify"
)

type AutomationPayload struct {
	NfcTagUid string `json:"nfc"`
	MediaId   string `json:"track"`
	Name      string `json:"name"`
}

func (r *Router) UseAutomation() {
	r.App.Get("/api/automations", func(c *fiber.Ctx) error {
		user := GetCurrentUser(c)
		if user.ID == 0 {
			return c.SendStatus(400)
		}

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

	r.App.Delete("/api/automation/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.SendStatus(400)
		}

		user := GetCurrentUser(c)
		if user.ID == 0 {
			return c.SendStatus(400)
		}

		automation := database.Automation{}
		if err := r.Conn.GetAutomation(uint(id), &automation); err != nil {
			if errors.Is(err, fiber.ErrNotFound) {
				return c.SendStatus(404)
			}
			return c.SendStatus(500)
		}

		if user.UserID != automation.UserID {
			return c.SendStatus(401)
		}

		if err := r.Conn.DeleteAutomation(uint(id)); err != nil {
			return c.SendStatus(500)
		}

		return c.SendStatus(200)
	})

	r.App.Post("/api/automation", func(c *fiber.Ctx) error {
		user := GetCurrentUser(c)
		if user.ID == 0 {
			return c.SendStatus(400)
		}

		payload := AutomationPayload{}

		if err := c.BodyParser(&payload); err != nil {
			log.Println("error = ", err)
			return err
		}

		log.Print(payload)

		nfcTag := database.NfcTag{}

		if err := r.Conn.GetTag(payload.NfcTagUid, &nfcTag); err != nil {
			log.Println("error = ", err)
			return err
		}

		automation := database.Automation{}

		if err := r.Conn.GetAutomationByNfcTag(nfcTag, &automation); err == nil {
			return c.SendStatus(409)
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Conn.CreateAutomation(nfcTag, payload.MediaId, payload.Name, user)
			return c.SendStatus(200)
		}

		return c.SendStatus(400)
	})

	r.App.Post("/api/automation/trigger/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.SendStatus(400)
		}

		user := GetCurrentUser(c)
		if user.ID == 0 {
			return c.SendStatus(400)
		}

		automation := database.Automation{}
		if err := r.Conn.GetAutomation(uint(id), &automation); err != nil {
			if errors.Is(err, fiber.ErrNotFound) {
				return c.SendStatus(404)
			}
			return c.SendStatus(500)
		}

		if user.UserID != automation.UserID {
			return c.SendStatus(401)
		}

		sessionId := c.Get("Session-Id")

		authData, err := spotify.GetAuthData(r.Conn, sessionId)
		if err != nil {
			return c.SendStatus(401)
		}

		err = spotify.StartPlayback(authData, sessionId, r.Conn, automation.MediaId)
		if err != nil {
			return c.SendStatus(500)
		}

		return c.SendStatus(200)
	})
}
