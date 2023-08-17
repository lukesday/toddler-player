package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"toddler-player/server/database"
)

func (r *Router) UseNfc() {
	r.App.Get("/api/nfc/unused", func(c *fiber.Ctx) error {
		nfcTags := []database.NfcTag{}
		if err := r.Conn.GetUnused(&nfcTags); err != nil {
			return err
		}

		return c.JSON(nfcTags)
	})

	r.App.Post("/api/nfc/:uid", func(c *fiber.Ctx) error {
		nfcTag := database.NfcTag{}
		// This can likely be refactored into a single query join
		if tagErr := r.Conn.GetTag(c.Params("uid"), &nfcTag); tagErr == nil {
			automation := database.Automation{}

			if err := r.Conn.GetAutomationByNfcTag(nfcTag, &automation); err != nil {
				log.Println(err)
			} else {
				// Trigger automation
				log.Println(automation)
			}
		} else if tagErr == gorm.ErrRecordNotFound {
			r.Conn.Enroll(c.Params("uid"))
		}

		return c.SendStatus(200)
	})
}
