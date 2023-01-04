package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbName     = "DB_NAME"
)

// ConnectDB conecta a la base de datos y devuelve una conexión abierta.
func ConnectDB() (*sql.DB, error) {
	host := os.Getenv(dbHost)
	port := os.Getenv(dbPort)
	user := os.Getenv(dbUser)
	password := os.Getenv(dbPassword)
	name := os.Getenv(dbName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
