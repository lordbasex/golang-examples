package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

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
	err = db.QueryRow("SELECT id, username FROM users WHERE username = ? AND password = ?", username, password).Scan(&user.ID, &user.User)
	if err != nil {
		return nil, err
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

func GetDataBase() (string, error) {

	var versionString string

	db, err := ConnectDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Realiza la consulta a la base de datos
	row := db.QueryRow("SELECT version()")

	if err := row.Scan(&versionString); err != nil {
		log.Print(err)
		return "", err
	}

	return versionString, nil
}

// SetCreditCard
func (c CreateCredit) SetCreditCard() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	sql := `
	INSERT INTO credit_card (uniqueid, card, expiration_date, security_code, amount, id_customer, created_at)
	VALUES (
		?,
		HEX(AES_ENCRYPT(?, ?)),
		HEX(AES_ENCRYPT(?, ?)),
		HEX(AES_ENCRYPT(?, ?)),
		HEX(AES_ENCRYPT(?, ?)),
		?,
		NOW()
	)`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		c.UniqueID,
		c.Card, AesKey,
		c.ExpirationDate, AesKey,
		c.SecurityCode, AesKey,
		c.Amount, AesKey,
		c.IDCustomer,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetCreditCardAll
func GetCreditCardAll() ([]CreateCredit, error) {

	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var getCreateCredit []CreateCredit

	sql := `
	SELECT
		id,
    	uniqueid,
    	AES_DECRYPT(UNHEX(card), ?) AS card,
    	AES_DECRYPT(UNHEX(expiration_date), ?) AS expiration_date,
    	AES_DECRYPT(UNHEX(security_code), ?) AS security_code,
    	AES_DECRYPT(UNHEX(amount), ?) AS amount,
    	id_customer,
    	created_at
	FROM
    	credit_card
	`
	rows, err := db.Query(sql, AesKey, AesKey, AesKey, AesKey)

	for rows.Next() {
		var data CreateCredit
		if err := rows.Scan(&data.ID, &data.UniqueID, &data.Card, &data.ExpirationDate, &data.SecurityCode, &data.Amount, &data.IDCustomer, &data.CreationAt); err != nil {
			return nil, err
		}
		getCreateCredit = append(getCreateCredit, data)
	}

	return getCreateCredit, err
}

// GetCreditCardByUniqueID
func GetCreditCardByUniqueID(uniqueID string) (CreateCredit, error) {

	var data CreateCredit

	db, err := ConnectDB()
	if err != nil {
		return data, err
	}

	sql := `
	SELECT
		id,
    	uniqueid,
    	AES_DECRYPT(UNHEX(card), ?) AS card,
    	AES_DECRYPT(UNHEX(expiration_date), ?) AS expiration_date,
    	AES_DECRYPT(UNHEX(security_code), ?) AS security_code,
    	AES_DECRYPT(UNHEX(amount), ?) AS amount,
    	id_customer,
    	created_at
	FROM
    	credit_card
	WHERE
		uniqueid = ?
	`

	row := db.QueryRow(sql, AesKey, AesKey, AesKey, AesKey, uniqueID)

	defer db.Close()

	if err := row.Scan(&data.ID, &data.UniqueID, &data.Card, &data.ExpirationDate, &data.SecurityCode, &data.Amount, &data.IDCustomer, &data.CreationAt); err != nil {
		if err != nil {
			return data, nil
		}
		return data, err
	}

	return data, err
}

// SetCreateSurvey
func (s Survey) SetCreateSurvey() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	sql := `INSERT INTO survey 
		(uniqueid, status, agent, queue, phone, id_customer, created_at)
	VALUES 
		(?, '0', ?, ?, ?, ?, now())`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		s.UniqueID,
		s.Agent,
		s.Queue,
		s.Phone,
		s.IDCustomer,
	)

	if err != nil {
		return err
	}

	return nil
}

func SetUpdateSurvey(uniqueID, q string, value int) error {
	//log.Printf("uniqueID: %s q: %s value: %d", uniqueID, q, value)

	db, err := ConnectDB()
	if err != nil {
		return err
	}

	var sql string
	switch q {
	case "q1":
		sql = fmt.Sprintf("UPDATE survey SET q1 = %d WHERE uniqueid = %s;", value, uniqueID)
	case "q2":
		sql = fmt.Sprintf("UPDATE survey SET q2 = %d WHERE uniqueid = %s;", value, uniqueID)
	case "q3":
		sql = fmt.Sprintf("UPDATE survey SET q3 = %d, status = 1 WHERE uniqueid = %s;", value, uniqueID)
	default:
		return err
	}

	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func SetAbortedSurvey(uniqueID string) error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Query the current status
	var status int
	query := "SELECT status FROM survey WHERE uniqueid = ?"
	err = db.QueryRow(query, uniqueID).Scan(&status)
	if err != nil {
		return err
	}

	if status == 0 {
		// Update status to 2
		updateQuery := "UPDATE survey SET status = 2 WHERE uniqueid = ?"
		_, err = db.Exec(updateQuery, uniqueID)
		if err != nil {
			return err
		}
	}

	return nil
}
