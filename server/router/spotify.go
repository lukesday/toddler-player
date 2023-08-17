package router

import (
	"log"
	"toddler-player/server/spotify"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) UseSpotify() {
	r.App.Get("/api/spotify/devices", func(c *fiber.Ctx) error {

		sessionId := c.Get("Session-Id")

		authData, err := spotify.GetAuthData(r.Conn, sessionId)
		if err != nil {
			return c.SendStatus(401)
		}

		if devices, deviceErr := spotify.GetDevices(authData, sessionId, r.Conn); deviceErr == nil {
			devices.SessionId = sessionId
			return c.JSON(devices)
		} else {
			log.Print(deviceErr)
		}

		return c.SendStatus(500)
	})

	r.App.Get("/api/spotify/me", func(c *fiber.Ctx) error {

		sessionId := c.Get("Session-Id")

		authData, err := spotify.GetAuthData(r.Conn, sessionId)
		if err != nil {
			return c.SendStatus(401)
		}

		if userData, err := spotify.GetUserData(authData, sessionId, r.Conn); err != nil {
			log.Print(err)
			return c.SendStatus(500)
		} else {
			userData.SessionId = sessionId
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

		if userData, err := spotify.GetUserDataWithoutRetry(authData); err != nil {
			return err
		} else {
			userId, _ = spotify.SaveAuthData(r.Conn, userId, userData.ID, authData)
			userData.SessionId = userId
			return c.JSON(userData)
		}
	})
}
