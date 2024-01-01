package handlers

import (
	"api-fiber/config"
	"api-fiber/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Hello is an HTTP request handler that responds with a protected message.
func Hello(c *fiber.Ctx) error {
	// Log the function if debug mode is enabled
	if config.DebugConfiguration {
		log.Print("Func handler.Hello")
	}

	// Get user information from the context
	claims, ok := c.Locals("user").(*models.CustomClaims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User information not available"})
	}

	// Protected response
	return c.JSON(fiber.Map{
		"message": "Hello, " + claims.UserInfo.FullName + "!",
	})
}
