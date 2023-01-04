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

	// Crea un nuevo enrutador.
	router := mux.NewRouter()

	// Agrega la ruta de inicio de sesión.
	router.HandleFunc("/login", loginHandler).Methods("POST")

	// Agrega otras rutas aquí...

	// Inicia el servidor.
	log.Fatal(http.ListenAndServe(":8080", router))
}
