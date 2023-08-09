package router

import (
	"github.com/gofiber/fiber/v2"
)

func (r *Router) UseDbLogs() {
	r.App.Get("/api/logs", func(c *fiber.Ctx) error {
		return c.JSON(r.Conn.GetRequestLogs())
	})
}
