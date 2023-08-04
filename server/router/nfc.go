package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"toddler-player/server/database"
)

func (r *Router) UseNfc() {
	r.App.Post("/api/nfc/:uid", func(c *fiber.Ctx) error {
		nfcTag := database.NfcTag{}
		if tagErr := r.Conn.GetTag(c.Params("uid"), nfcTag); tagErr == nil {
			automation := database.Automation{}

			if err := r.Conn.GetAutomation(nfcTag, automation); err != nil {
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
