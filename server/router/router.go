package router

import (
	"github.com/gofiber/fiber/v2"

	"toddler-player/server/database"
)

type Router struct {
	App  *fiber.App
	Conn *database.DatabaseConnection
}

func NewRouter(conn *database.DatabaseConnection) *Router {
	router := Router{
		App:  fiber.New(),
		Conn: conn,
	}
	return &router
}

func InitialiseRouter(conn *database.DatabaseConnection) {
	router := NewRouter(conn)

	router.UseMiddleware()
	router.UseDbLogs()
	router.UseDevice()
	router.UseNfc()
	router.UseAutomation()
	router.UseSpotify()

	go router.App.Listen(":3000")
}
