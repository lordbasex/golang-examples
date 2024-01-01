package router

import (
	"api-fiber/handlers"

	"github.com/gofiber/fiber/v2"
)

func LoginRoutes(app *fiber.App) {
	app.Post("/login", handlers.Login)
}
