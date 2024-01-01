package router

import (
	"api-fiber/handlers"

	"github.com/gofiber/fiber/v2"
)

func HelloRoutes(app *fiber.App) {
	app.Get("/api/v1/hello", handlers.Hello)
}
