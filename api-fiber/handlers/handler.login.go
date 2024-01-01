package handlers

import (
	"api-fiber/config"
	"api-fiber/middleware"
	"api-fiber/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Login is an HTTP request handler for user authentication.
func Login(c *fiber.Ctx) error {
	// Log the function if debug mode is enabled
	if config.DebugConfiguration {
		log.Print("Func handler.Login")
	}

	// Hardcoded user (for demonstration purposes)
	user := models.User{FullName: "Federico Pereira", Email: "lord.basex@gmail.com", Password: "I'llBeBack"}

	// Validate credentials
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check credentials against hardcoded values
	if user.Email != "lord.basex@gmail.com" || user.Password != "I'llBeBack" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Create a token
	token, err := middleware.CreateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token creation failed"})
	}

	// Successful login response
	return c.JSON(fiber.Map{
		"status": 200,
		"token":  token,
		"msg":    "Successful login",
	})
}
