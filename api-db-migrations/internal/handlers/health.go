package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// HealthHandler maneja el endpoint de health check
// @Summary Health check
// @Description Verifica que el servidor esté funcionando
// @Tags health
// @Success 200 {string} string "ok"
// @Router /health [get]
func HealthHandler(c *fiber.Ctx) error {
	return c.SendString("ok")
}

