package router

import (
	"api-fiber/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// FiberMiddleware sets up middleware for a Fiber web application.
func FiberMiddleware(app *fiber.App) {
	// Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	// Middleware to validate the token on routes /api/v1/...
	app.Use("/api/v1/*", middleware.ValidateToken)

	// List of blocked IPs
	blockedIPs := []string{"172.16.0.33", "10.0.0.2"}

	// IP blocking middleware
	app.Use(middleware.BlockIP(blockedIPs))
}
