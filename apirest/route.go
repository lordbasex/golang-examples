package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() {
	router := mux.NewRouter()

	router.Use(jwtMiddleware)

	router.HandleFunc("/", root).Methods("GET")

	router.HandleFunc("/login", authMiddleware(loginHandler)).Methods("POST")

	router.HandleFunc("/api/v1/hello", authMiddleware(helloWorld)).Methods("GET")

	router.HandleFunc("/api/v1/customer/document_number/{id:[0-9]+}", authMiddleware(customerByDocumentNumberHandler)).Methods("GET")
	router.HandleFunc("/api/v1/customer/customer_number/{id:[0-9]+}", authMiddleware(customerByCustomerNumberHandler)).Methods("GET")

	router.HandleFunc("/api/v1/creditcard/create", authMiddleware(CreateCreditCardHandler)).Methods("POST")
	router.HandleFunc("/api/v1/creditcard/search", authMiddleware(CreateCreditAll)).Methods("GET")
	router.HandleFunc("/api/v1/creditcard/search/{uniqueid:[0-9.]+}", authMiddleware(CreateCreditSearch)).Methods("GET")

	router.HandleFunc("/api/v1/survey/create", authMiddleware(CreateSurvey)).Methods("POST")
	router.HandleFunc("/api/v1/survey/update/{uniqueid:[0-9.]+}", authMiddleware(UpdateSurvey)).Methods("PUT")
	router.HandleFunc("/api/v1/survey/aborted/{uniqueid:[0-9.]+}", authMiddleware(AbortedSurvey)).Methods("PUT")

	//Start Service
	log.Println("RESTAPI: Start Port: 1444")

	// Inicia el servidor.
	log.Fatal(http.ListenAndServe(":1444", router))
}
