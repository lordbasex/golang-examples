package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"api/internal/handlers"
)

// Setup configura todas las rutas de la aplicación
func Setup(app *fiber.App, db *sql.DB) {
	// Middlewares globales
	app.Use(recover.New()) // Recuperación de panics
	app.Use(logger.New())  // Logging de requests

	// Inicializar handlers
	dbHandler := handlers.NewDBHandler(db)
	userHandler := handlers.NewUserHandler(db)
	campaignHandler := handlers.NewCampaignHandler(db)

	// Rutas públicas
	app.Get("/health", handlers.HealthHandler)
	app.Get("/db", dbHandler.ShowTables)

	// Grupo API v1
	api := app.Group("/api/v1")

	// Rutas de usuarios
	users := api.Group("/users")
	users.Get("/", userHandler.GetAll)
	users.Get("/:id", userHandler.GetByID)
	users.Post("/", userHandler.Create)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	// Rutas de campañas
	campaigns := api.Group("/campaigns")
	campaigns.Get("/stats", campaignHandler.GetStats) // Debe ir antes de /:id
	campaigns.Get("/", campaignHandler.GetAll)
	campaigns.Get("/:id", campaignHandler.GetByID)
	campaigns.Post("/", campaignHandler.Create)
	campaigns.Put("/:id", campaignHandler.Update)
	campaigns.Delete("/:id", campaignHandler.Delete)
}
