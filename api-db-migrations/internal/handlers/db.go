package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// DBHandler contiene las dependencias para los handlers de base de datos
type DBHandler struct {
	DB *sql.DB
}

// NewDBHandler crea una nueva instancia de DBHandler
func NewDBHandler(db *sql.DB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}

// ShowTables muestra todas las tablas de la base de datos
// @Summary Listar tablas
// @Description Obtiene la lista de todas las tablas en la base de datos
// @Tags database
// @Produce plain
// @Success 200 {string} string "Lista de tablas"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /db [get]
func (h *DBHandler) ShowTables(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SHOW TABLES")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		tables = append(tables, name)
	}
	
	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Devolver las tablas como texto plano
	result := ""
	for _, t := range tables {
		result += t + "\n"
	}
	
	return c.SendString(result)
}

