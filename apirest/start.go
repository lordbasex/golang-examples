package main

import "log"

var DebugConfiguration bool
var Secret []byte
var AesKey string

func Start() {
	DebugConfiguration = false
	Secret = []byte("my-secret")
	AesKey = "demo"

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
