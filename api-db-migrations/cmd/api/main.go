package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"api/internal/db"
	"api/internal/router"
)

func main() {
	// Obtener configuración desde variables de entorno
	migrationsPath := getenv("MIGRATIONS_PATH", "file://migrations")
	addr := getenv("HTTP_ADDR", ":8080")

	// Conectar a la base de datos
	sqlDB := db.MustConnectDB()
	defer sqlDB.Close()

	// 1) migraciones al iniciar (create/alter según versión)
	db.MustRunMigrations(sqlDB, migrationsPath)

	// 2) configurar Fiber
	app := fiber.New(fiber.Config{
		AppName:      "API DB Migrations",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// 3) configurar rutas y middlewares
	router.Setup(app, sqlDB)

	// 4) iniciar servidor
	log.Printf("api: listening on %s", addr)
	log.Fatal(app.Listen(addr))
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
