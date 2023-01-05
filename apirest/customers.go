package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Customer represents a customer in the database
type Customer struct {
	ID             int
	Name           string
	LastName       string
	DocumentNumber int
	CustomerNumber int
	Address        string
	Phone          string
	Email          string
	Enabled        bool
	CreatedAt      string
	UpdatedAt      string
	LastEntryIVR   string
}

// customersHandler
func customersHandler(w http.ResponseWriter, r *http.Request) {

	// Obtiene los parámetros de búsqueda de la solicitud.
	queryParams := r.URL.Query()
	documentNumberStr := queryParams.Get("document_number")
	customerNumberStr := queryParams.Get("customer_number")

	if documentNumberStr != "" {
		// Filtra por número de documento
		documentNumber, err := strconv.Atoi(documentNumberStr)
		if err != nil {
			http.Error(w, "Parámetro de número de documento inválido", http.StatusBadRequest)
			return
		}

		// Obtiene el cliente de la base de datos
		customer, err := GetCustomerByDocumentNumber(documentNumber)
		if err != nil {
			if err.Error() == "cliente no encontrado" {
				http.Error(w, "Cliente no encontrado", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Devuelve el cliente en formato JSON
		customerJSON, err := json.Marshal(customer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(customerJSON)
	} else if customerNumberStr != "" {
		// Filtra por número de cliente
		customerNumber, err := strconv.Atoi(customerNumberStr)
		if err != nil {
			http.Error(w, "Parámetro de número de cliente inválido", http.StatusBadRequest)
			return
		}

		// Obtiene el cliente de la base de datos
		customer, err := GetCustomerByCustomerNumber(customerNumber)
		if err != nil {
			if err.Error() == "cliente no encontrado" {
				http.Error(w, "Cliente no encontrado", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Devuelve el cliente en formato JSON
		customerJSON, err := json.Marshal(customer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(customerJSON)
	}
}

// customersListHandler maneja las solicitudes GET al endpoint /customers/list.
// Verifica el token JWT y, si es válido, obtiene todos los clientes de la base de datos y los devuelve en formato JSON.
func customersListHandler(w http.ResponseWriter, r *http.Request) {
	// Obtiene los clientes de la base de datos
	customers, err := GetCustomers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devuelve los clientes en formato JSON
	customersJSON, err := json.Marshal(customers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(customersJSON)
}
