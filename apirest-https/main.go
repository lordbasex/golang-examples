package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Msg struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/health", health).Methods("GET")

	// For dev only - Set up CORS so React client can consume our API
	corsWrapper := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Content-Type", "Origin", "Accept", "*"},
		AllowCredentials: true,
		Debug:            false,
	})

	//Start Service
	log.Println("RESTAPI: Start")

	err := http.ListenAndServeTLS(":443", "localhost.crt", "localhost.key", corsWrapper.Handler(router))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Msg{
		Status: http.StatusOK,
		Msg:    "OK",
	})
}
