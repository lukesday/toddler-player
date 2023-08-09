package router

import (
	"log"
	"toddler-player/server/spotify"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) UseSpotify() {
	r.App.Get("/api/spotify/devices", func(c *fiber.Ctx) error {

		authData, err := spotify.GetAuthData(r.Conn, c.Get("Session-Id"))
		if err != nil {
			return c.SendStatus(401)
		}

		if devices, err := spotify.GetDevices(authData, c.Get("Session-Id"), r.Conn); err == nil {
			return c.JSON(devices)
		}

		return c.SendStatus(401)
	})

	r.App.Get("/api/spotify/me", func(c *fiber.Ctx) error {

		authData, err := spotify.GetAuthData(r.Conn, c.Get("Session-Id"))
		if err != nil {
			return c.SendStatus(401)
		}

		if userData, err := spotify.GetUserData(authData, c.Get("Session-Id"), r.Conn); err != nil {
			log.Print(err)
			return c.SendStatus(500)
		} else {
			return c.JSON(userData)
		}
	})

	r.App.Post("/api/spotify/auth", func(c *fiber.Ctx) error {
		payload := spotify.SpotifyAuthPayload{}

		if err := c.BodyParser(&payload); err != nil {
			log.Println("error = ", err)
			return err
		}

		var authData spotify.SpotifyAuthResponse

		if tokenData, err := spotify.GetAccessToken(payload); err != nil {
			return c.SendStatus(401)
		} else {
			authData = tokenData
		}

		userId := c.Get("Session-Id")
		spotify.SaveAuthData(r.Conn, userId, authData)

		if userData, err := spotify.GetUserData(authData, c.Get("Session-Id"), r.Conn); err != nil {
			return err
		} else {
			userData.SessionId = userId
			return c.JSON(userData)
		}
	})
}
