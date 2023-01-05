package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

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
			(1, "user1", "pass1"),
			(2, "user2", "pass2"),
			(3, "user3", "pass3"),
			(4, "user4", "pass4"),
			(5, "user5", "pass5"),
			(6, "user6", "pass6"),
			(7, "user7", "pass7"),
			(8, "user8", "pass8"),
			(9, "user9", "pass9"),
			(10, "user10", "pass10")
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
		(1, 'Juan', 'Pérez', 1234567890, 1001, 'Calle Falsa 123', '11-1111-1111', 'juan.perez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(2, 'Ana', 'Rodríguez', 1234567891, 1002, 'Calle Falsa 456', '11-2222-2222', 'ana.rodriguez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(3, 'Mario', 'González', 1234567892, 1003, 'Calle Falsa 789', '11-3333-3333', 'mario.gonzalez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(4, 'Lola', 'Gómez', 1234567893, 1004, 'Calle Falsa 321', '11-4444-4444', 'lola.gomez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(5, 'Jorge', 'Sánchez', 1234567894, 1005, 'Calle Falsa 159', '11-5555-5555', 'jorge.sanchez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(6, 'Sofía', 'Martínez', 1234567895, 1006, 'Calle Falsa 753', '11-6666-6666', 'sofia.martinez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(7, 'Pablo', 'Díaz', 1234567896, 1007, 'Calle Falsa 951', '11-7777-7777', 'pablo.diaz@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(8, 'Rocío', 'Castro', 1234567897, 1008, 'Calle Falsa 147', '11-8888-8888', 'rocio.castro@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(9, 'Julián', 'Suárez', 1234567898, 1009, 'Calle Falsa 753', '11-9999-9999', 'julian.suarez@ejemplo.com', '1', NOW(), NOW(), NOW()),
		(10, 'Evelyn', 'Rodríguez', 1234567899, 1010, 'Calle Falsa 159', '11-0000-0000', 'evelyn.rodriguez@ejemplo.com', '1', NOW(), NOW(), NOW())
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

// GetCustomer devuelve el cliente con el número de documento especificado de la base de datos
func GetCustomer(documentNumber int) (*Customer, error) {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Ejecuta una consulta SELECT para recuperar el cliente con el número de documento especificado
	row := db.QueryRow("SELECT * FROM customers WHERE document_number = ?", documentNumber)

	// Escanea el resultado en una variable
	var id int
	var name, lastName, address, phone, email string
	var customerNumber int
	var enabled bool
	var createdAt, updatedAt, lastEntryIvr string
	err = row.Scan(&id, &name, &lastName, &documentNumber, &customerNumber, &address, &phone, &email, &enabled, &createdAt, &updatedAt, &lastEntryIvr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente no encontrado")
		}
		return nil, err
	}

	// Devuelve los datos del cliente
	return &Customer{
		ID:             id,
		Name:           name,
		LastName:       lastName,
		DocumentNumber: documentNumber,
		CustomerNumber: customerNumber,
		Address:        address,
		Phone:          phone,
		Email:          email,
		Enabled:        enabled,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		LastEntryIVR:   lastEntryIvr,
	}, nil
}

// GetCustomerByDocumentNumber obtiene un cliente de la base de datos según su número de documento.
func GetCustomerByDocumentNumber(documentNumber int) (Customer, error) {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return Customer{}, err
	}
	defer db.Close()

	// Realiza la consulta a la base de datos
	row := db.QueryRow("SELECT * FROM customers WHERE document_number = ?", documentNumber)

	// Obtiene el cliente de la fila
	var customer Customer
	err = row.Scan(&customer.ID, &customer.Name, &customer.LastName, &customer.DocumentNumber, &customer.CustomerNumber, &customer.Address, &customer.Phone, &customer.Email, &customer.Enabled, &customer.CreatedAt, &customer.UpdatedAt, &customer.LastEntryIVR)
	if err == sql.ErrNoRows {
		return Customer{}, errors.New("cliente no encontrado")
	} else if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

// GetCustomerByCustomerNumber obtiene un cliente de la base de datos según su número de cliente.
func GetCustomerByCustomerNumber(customerNumber int) (Customer, error) {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return Customer{}, err
	}
	defer db.Close()

	// Realiza la consulta a la base de datos
	row := db.QueryRow("SELECT * FROM customers WHERE customer_number = ?", customerNumber)

	// Obtiene el cliente de la fila
	var customer Customer
	err = row.Scan(&customer.ID, &customer.Name, &customer.LastName, &customer.DocumentNumber, &customer.CustomerNumber, &customer.Address, &customer.Phone, &customer.Email, &customer.Enabled, &customer.CreatedAt, &customer.UpdatedAt, &customer.LastEntryIVR)
	if err == sql.ErrNoRows {
		return Customer{}, errors.New("cliente no encontrado")
	} else if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

// GetCustomers obtiene todos los clientes de la base de datos.
func GetCustomers() ([]Customer, error) {
	// Conecta a la base de datos.
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Realiza la consulta a la base de datos
	rows, err := db.Query("SELECT * FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Obtiene los clientes de las filas
	var customers []Customer
	for rows.Next() {
		var customer Customer
		err = rows.Scan(&customer.ID, &customer.Name, &customer.LastName, &customer.DocumentNumber, &customer.CustomerNumber, &customer.Address, &customer.Phone, &customer.Email, &customer.Enabled, &customer.CreatedAt, &customer.UpdatedAt, &customer.LastEntryIVR)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}
