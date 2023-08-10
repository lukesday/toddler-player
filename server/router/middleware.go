package router

import (
	"toddler-player/server/database"

	"github.com/gofiber/fiber/v2"
)

func (r *Router) UseMiddleware() {
	r.App.All("/api/*", func(c *fiber.Ctx) error {
		r.Conn.LogRequest(c.Request().URI().String(), string(c.Request().Header.Method()))
		return c.Next()
	})

	r.App.All("/api/*", func(c *fiber.Ctx) error {
		sessionId := c.Get("Session-Id")
		if sessionId == "" {
			return c.Next()
		}

		var user database.User
		if err := r.Conn.GetUser(sessionId, &user); err == nil {
			c.Locals("user", user)
		}

		return c.Next()
	})
}

func GetCurrentUser(c *fiber.Ctx) database.User {
	var user database.User
	localUser := c.Locals("user")
	if localUser == nil {
		return user
	}
	user = localUser.(database.User)

	return user
}
