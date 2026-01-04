package models

import "time"

// User representa un usuario en la base de datos
type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     *string   `json:"phone,omitempty"` // Pointer para permitir NULL
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	FullName string  `json:"full_name" validate:"required,min=2,max=255"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=32"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Email    string  `json:"email,omitempty" validate:"omitempty,email"`
	FullName string  `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=32"`
}

// UserResponse representa la respuesta con información del usuario
type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    *User  `json:"data,omitempty"`
}

// UsersResponse representa la respuesta con lista de usuarios
type UsersResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
	Data    []*User `json:"data,omitempty"`
	Total   int     `json:"total"`
}

