package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

// User representa una cuenta de usuario.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// createUsersTable crea la tabla de usuarios en la base de datos.
func createUsersTable() error {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Crea la tabla de usuarios.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Inserta ejemplos de usuarios.
	_, err = db.Exec(`
		INSERT INTO users (username, password) VALUES
			("user1", "pass1"),
			("user2", "pass2"),
			("user3", "pass3"),
			("user4", "pass4"),
			("user5", "pass5"),
			("user6", "pass6"),
			("user7", "pass7"),
			("user8", "pass8"),
			("user9", "pass9"),
			("user10", "pass10")
	`)
	return err
}

// verifyCredentials verifica si las credenciales de inicio de sesión proporcionadas son válidas.
func verifyCredentials(username, password string) (*User, error) {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Verifica si las credenciales de inicio de sesión proporcionadas son válidas.
	var user User
	err = db.QueryRow("SELECT id, username FROM users WHERE username = ? AND password = ?", username, password).Scan(&user.ID, &user.Username)
	if err == sql.ErrNoRows {
		return nil, errors.New("Credenciales de inicio de sesión no válidas")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
