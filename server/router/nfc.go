package router

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"toddler-player/server/database"
	"toddler-player/server/spotify"
)

func (r *Router) UseNfc() {
	r.App.Get("/api/nfc/unused", func(c *fiber.Ctx) error {
		nfcTags := []database.NfcTag{}
		if err := r.Conn.GetUnusedTag(&nfcTags); err != nil {
			return err
		}

		return c.JSON(nfcTags)
	})

	r.App.Post("/api/nfc/:uid", func(c *fiber.Ctx) error {
		nfcTag := database.NfcTag{}
		if tagErr := r.Conn.GetTag(c.Params("uid"), &nfcTag); tagErr == nil {
			automation := database.Automation{}

			if err := r.Conn.GetAutomationByNfcTag(nfcTag, &automation); err != nil {
			} else {
				userId := automation.UserID

				fmt.Println("userId", userId)

				authData, err := spotify.GetAuthData(r.Conn, userId)
				if err != nil {
					fmt.Println("error getting auth data", err)
					return c.SendStatus(401)
				}
				// Trigger automation
				err = spotify.StartPlayback(authData, userId, r.Conn, automation.MediaId)
				if err != nil {
					fmt.Println("error starting playback", err)
					return c.SendStatus(500)
				}
			}
		} else if tagErr == gorm.ErrRecordNotFound {
			r.Conn.Enroll(c.Params("uid"))
		}

		return c.SendStatus(200)
	})
}
