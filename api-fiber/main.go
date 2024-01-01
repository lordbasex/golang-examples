package main

import (
	"api-fiber/config"
	"api-fiber/database"
	"api-fiber/middleware"
	"api-fiber/router"
	"api-fiber/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	// DebugAPI
	config.DebugAPI()

	// Generate RSA keys
	middleware.PrivateKey, middleware.PublicKey = utils.GenarateRSAKeys()

	// create tables
	database.DBInit()
}

func main() {
	// Define Fiber config.
	c := config.FiberConfig()

	// Define a new Fiber app with config
	app := fiber.New(c)

	// Middleware
	router.FiberMiddleware(app)

	// Logger
	app.Use(logger.New())

	// Router
	router.LoginRoutes(app)   //POST - /login
	router.HelloRoutes(app)   //GET - /api/v1/hello
	router.NotFoundRoute(app) // NotFoundRoute

	// Start server
	utils.StartServer(app)
}
