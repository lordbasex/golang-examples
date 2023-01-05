package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Crea la tabla de usuarios en la base de datos y inserta ejemplos de usuarios.
	if err := createUsersTable(); err != nil {
		log.Fatal(err)
	}

	// Crea la tabla de customer en la base de datos y inserta ejemplos de customer.
	if err := createCustomersTable(); err != nil {
		log.Fatal(err)
	}

	// Crea un nuevo enrutador.
	router := mux.NewRouter()

	// Pasa todo por el Middleware
	router.Use(jwtMiddleware)

	// Agrega la ruta de inicio de sesión.
	router.HandleFunc("/login", authMiddleware(loginHandler)).Methods("POST")
	router.HandleFunc("/customers", authMiddleware(customersHandler)).Methods("GET")
	router.HandleFunc("/customers/list", authMiddleware(customersListHandler)).Methods("GET")

	// Inicia el servidor.
	log.Fatal(http.ListenAndServe(":8080", router))
}
