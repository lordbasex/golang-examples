package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"api/internal/models"
)

// UserHandler contiene las dependencias para los handlers de usuarios
type UserHandler struct {
	DB *sql.DB
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		DB: db,
	}
}

// GetAll obtiene todos los usuarios
// @Summary Listar usuarios
// @Description Obtiene la lista de todos los usuarios
// @Tags users
// @Produce json
// @Success 200 {object} models.UsersResponse
// @Failure 500 {object} models.UsersResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	query := `SELECT id, email, full_name, phone, created_at FROM users ORDER BY created_at DESC`

	rows, err := h.DB.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.UsersResponse{
			Success: false,
			Message: "Error al obtener usuarios: " + err.Error(),
		})
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.FullName, &user.Phone, &user.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.UsersResponse{
				Success: false,
				Message: "Error al leer usuarios: " + err.Error(),
			})
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.UsersResponse{
			Success: false,
			Message: "Error al procesar usuarios: " + err.Error(),
		})
	}

	return c.JSON(models.UsersResponse{
		Success: true,
		Data:    users,
		Total:   len(users),
	})
}

// GetByID obtiene un usuario por su ID
// @Summary Obtener usuario por ID
// @Description Obtiene un usuario específico por su ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.UserResponse
// @Failure 404 {object} models.UserResponse
// @Failure 500 {object} models.UserResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.UserResponse{
			Success: false,
			Message: "ID de usuario inválido",
		})
	}

	query := `SELECT id, email, full_name, phone, created_at FROM users WHERE id = ?`

	user := &models.User{}
	err = h.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Phone,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(models.UserResponse{
			Success: false,
			Message: "Usuario no encontrado",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.UserResponse{
			Success: false,
			Message: "Error al obtener usuario: " + err.Error(),
		})
	}

	return c.JSON(models.UserResponse{
		Success: true,
		Data:    user,
	})
}

// Create crea un nuevo usuario
// @Summary Crear usuario
// @Description Crea un nuevo usuario en la base de datos
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.UserResponse
// @Failure 409 {object} models.UserResponse
// @Failure 500 {object} models.UserResponse
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	req := new(models.CreateUserRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.UserResponse{
			Success: false,
			Message: "Datos inválidos: " + err.Error(),
		})
	}

	// Validación básica
	if req.Email == "" || req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.UserResponse{
			Success: false,
			Message: "Email y nombre completo son requeridos",
		})
	}

	query := `INSERT INTO users (email, full_name, phone) VALUES (?, ?, ?)`

	result, err := h.DB.Exec(query, req.Email, req.FullName, req.Phone)
	if err != nil {
		// Check for duplicate email (unique constraint)
		if err.Error() == "Error 1062: Duplicate entry" {
			return c.Status(fiber.StatusConflict).JSON(models.UserResponse{
				Success: false,
				Message: "El email ya está registrado",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.UserResponse{
			Success: false,
			Message: "Error al crear usuario: " + err.Error(),
		})
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.UserResponse{
			Success: false,
			Message: "Error al obtener ID del usuario creado",
		})
	}

	// Obtener el usuario creado
	user := &models.User{}
	err = h.DB.QueryRow(
		`SELECT id, email, full_name, phone, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Phone, &user.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(models.UserResponse{
			Success: true,
			Message: "Usuario creado exitosamente",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.UserResponse{
		Success: true,
		Message: "Usuario creado exitosamente",
		Data:    user,
	})
}

// Update actualiza un usuario existente
// @Summary Actualizar usuario
// @Description Actualiza los datos de un usuario existente
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "Datos a actualizar"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.UserResponse
// @Failure 404 {object} models.UserResponse
// @Failure 500 {object} models.UserResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.UserResponse{
			Success: false,
			Message: "ID de usuario inválido",
		})
	}

	req := new(models.UpdateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.UserResponse{
			Success: false,
			Message: "Datos inválidos: " + err.Error(),
		})
	}

	// Verificar que el usuario existe
	var exists bool
	err = h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, userID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(models.UserResponse{
			Success: false,
			Message: "Usuario no encontrado",
		})
	}

	// Construir query dinámico solo con los campos proporcionados
	query := `UPDATE users SET `
	args := make([]interface{}, 0)
	updates := make([]string, 0)

	if req.Email != "" {
		updates = append(updates, "email = ?")
		args = append(args, req.Email)
	}
	if req.FullName != "" {
		updates = append(updates, "full_name = ?")
		args = append(args, req.FullName)
	}
	if req.Phone != nil {
		updates = append(updates, "phone = ?")
		args = append(args, req.Phone)
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.UserResponse{
			Success: false,
			Message: "No hay datos para actualizar",
		})
	}

	// Unir los updates y agregar la condición WHERE
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"
	args = append(args, userID)

	_, err = h.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.UserResponse{
			Success: false,
			Message: "Error al actualizar usuario: " + err.Error(),
		})
	}

	// Obtener el usuario actualizado
	user := &models.User{}
	err = h.DB.QueryRow(
		`SELECT id, email, full_name, phone, created_at FROM users WHERE id = ?`,
		userID,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Phone, &user.CreatedAt)

	if err != nil {
		return c.JSON(models.UserResponse{
			Success: true,
			Message: "Usuario actualizado exitosamente",
		})
	}

	return c.JSON(models.UserResponse{
		Success: true,
		Message: "Usuario actualizado exitosamente",
		Data:    user,
	})
}

// Delete elimina un usuario
// @Summary Eliminar usuario
// @Description Elimina un usuario de la base de datos
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.UserResponse
// @Failure 404 {object} models.UserResponse
// @Failure 500 {object} models.UserResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.UserResponse{
			Success: false,
			Message: "ID de usuario inválido",
		})
	}

	// Verificar que el usuario existe antes de eliminar
	var exists bool
	err = h.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, userID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(models.UserResponse{
			Success: false,
			Message: "Usuario no encontrado",
		})
	}

	query := `DELETE FROM users WHERE id = ?`
	_, err = h.DB.Exec(query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.UserResponse{
			Success: false,
			Message: "Error al eliminar usuario: " + err.Error(),
		})
	}

	return c.JSON(models.UserResponse{
		Success: true,
		Message: "Usuario eliminado exitosamente",
	})
}
