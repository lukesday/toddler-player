package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) UseMiddleware() {
	r.App.All("/api/*", func(c *fiber.Ctx) error {
		r.Conn.LogRequest(c.Request().URI().String(), string(c.Request().Header.Method()))
		return c.Next()
	})
}
