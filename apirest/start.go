package main

import (
	"log"
	"os"
)

const (
	debug       = "DEBUG"
	tokenSecret = "TOKEN_SECRET"
	dbEnc       = "MYSQL_ENC"
	dbHost      = "DB_HOST"
	dbPort      = "DB_PORT"
	dbUser      = "DB_USER"
	dbPassword  = "DB_PASSWORD"
	dbName      = "DB_NAME"
)

var DebugConfiguration bool
var Secret []byte
var AesKey string

func Start() {

	DebugConfiguration, _ = stringToBool(os.Getenv(debug))
	Secret = []byte(os.Getenv(tokenSecret))
	AesKey = os.Getenv(dbEnc)

	// Crea la tabla de usuarios en la base de datos y inserta ejemplos de usuarios.
	if err := createUsersTable(); err != nil {
		log.Fatal(err)
	}

	// Crea la tabla de customer en la base de datos y inserta ejemplos de customer.
	if err := createCustomersTable(); err != nil {
		log.Fatal(err)
	}

	// Crea la tabla de credit_card en la base de datos.
	if err := createCreditCardTable(); err != nil {
		log.Fatal(err)
	}

	// Crea la tabla de survey en la base de datos.
	if err := createSurveyTable(); err != nil {
		log.Fatal(err)
	}

}
