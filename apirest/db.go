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
		INSERT IGNORE users (id, username, password) VALUES
			(1, "superadmin", "admin"),
			(2, "user01", "pass1"),
			(3, "user02", "pass2"),
			(4, "user03", "pass3")
	`)
	return err
}

// createCustomersTable crea la tabla de usuarios en la base de datos.
func createCustomersTable() error {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Crea la tabla de customers.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id int(11) NOT NULL AUTO_INCREMENT,
			name varchar(100) NOT NULL,
			last_name varchar(100) NOT NULL,
			document_number int(11) NOT NULL,
			customer_number int(11) NOT NULL,
			address varchar(255) NOT NULL,
			phone varchar(30) NOT NULL DEFAULT '',
			email varchar(200) NOT NULL DEFAULT '',
			enabled varchar(2) NOT NULL DEFAULT '1',
			created_at timestamp NULL DEFAULT NULL,
			updated_at timestamp NULL DEFAULT NULL,
			last_entry_ivr timestamp NULL DEFAULT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY document_number (document_number),
			UNIQUE KEY customer_number (customer_number)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8
	`)
	if err != nil {
		return err
	}

	// Inserta ejemplos de customers.
	_, err = db.Exec(`
	INSERT IGNORE customers (id, name, last_name, document_number, customer_number, address, phone, email, enabled, created_at, updated_at, last_entry_ivr) VALUES
		(1, 'Federico', 'Pereira', 12345678, 1234567890, 'Calle Falsa 123', '11-1111-1111', 'juan.perez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(2, 'Ana', 'Rodríguez', 12345601, 1234567810, 'Calle Falsa 456', '11-2222-2222', 'ana.rodriguez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(3, 'Mario', 'González', 12345602, 1234567812, 'Calle Falsa 789', '11-3333-3333', 'mario.gonzalez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(4, 'Lola', 'Gómez', 12345603, 1234567813, 'Calle Falsa 321', '11-4444-4444', 'lola.gomez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(5, 'Jorge', 'Sánchez', 12345604, 1234567814, 'Calle Falsa 159', '11-5555-5555', 'jorge.sanchez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(6, 'Sofía', 'Martínez', 12345605, 1234567815, 'Calle Falsa 753', '11-6666-6666', 'sofia.martinez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(7, 'Pablo', 'Díaz', 12345606, 1234567816, 'Calle Falsa 951', '11-7777-7777', 'pablo.diaz@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(8, 'Rocío', 'Castro', 12345607, 1234567817, 'Calle Falsa 147', '11-8888-8888', 'rocio.castro@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(9, 'Julián', 'Suárez', 12345608, 1234567818, 'Calle Falsa 753', '11-9999-9999', 'julian.suarez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(10, 'Evelyn', 'Rodríguez', 12345609, 1234567819, 'Calle Falsa 159', '11-0000-0000', 'evelyn.rodriguez@ejemplo.com', '1', NOW(), NOW(), NOW())
	`)
	return err
}

// createCreditCardTable crea la tabla de credit_card en la base de datos.
func createCreditCardTable() error {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Crea la tabla de usuarios.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS credit_card (
			id int(11) NOT NULL AUTO_INCREMENT,
			uniqueid char(32) DEFAULT NULL,
			card blob,
			expiration_date blob,
			security_code blob,
			amount blob,
			id_customer int(11) DEFAULT NULL,
			created_at datetime DEFAULT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY uniqueid (uniqueid)
		)
	`)
	if err != nil {
		return err
	}

	return err
}

// createSurveyTable crea la tabla de credit_card en la base de datos.
func createSurveyTable() error {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Crea la tabla de survey.
	_, err = db.Exec(`

	CREATE TABLE IF NOT EXISTS survey (
			id int(10) NOT NULL AUTO_INCREMENT,
			uniqueid varchar(32) NOT NULL DEFAULT '',
			status tinyint(1) DEFAULT NULL,
			agent varchar(255) DEFAULT NULL,
			queue varchar(255) DEFAULT NULL,
			phone varchar(255) DEFAULT NULL,
			id_customer int(11) NOT NULL,
			created_at datetime DEFAULT NULL,
			q1 varchar(255) DEFAULT '',
			q2 varchar(255) DEFAULT '',
			q3 varchar(255) DEFAULT '',
			PRIMARY KEY (id),
			UNIQUE KEY uniqueid (uniqueid)
	)`)

	if err != nil {
		return err
	}

	return err
}
