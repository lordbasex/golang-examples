package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

// customerByDocumentNumberHandler
func customerByDocumentNumberHandler(w http.ResponseWriter, r *http.Request) {
	if DebugConfiguration {
		log.Print("func: customerByDocumentNumberHandler")
	}

	url := strings.Split(r.URL.Path, "/")

	documentNumber, err := strconv.Atoi(url[5])
	if err != nil || documentNumber < 1 {

		response := &ResponseCustomerERROR{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "wrong_parameter",
		}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	if DebugConfiguration {
		log.Printf("Document Number: %d", documentNumber)
	}

	// Obtiene el cliente de la base de datos
	customer, err := GetCustomerByDocumentNumber(documentNumber)
	if err != nil {
		response := &ResponseCustomerERROR{
			Code:    http.StatusNotFound,
			Status:  "success",
			Message: "not_exist",
		}
		responseWriter(http.StatusNotFound, response, w)
		return
	}

	response := &ResponseCustomerOK{
		Code:    http.StatusOK,
		Status:  "success",
		Message: customer,
	}
	responseWriter(http.StatusOK, response, w)
}

// customerByCustomerNumberHandler
func customerByCustomerNumberHandler(w http.ResponseWriter, r *http.Request) {
	if DebugConfiguration {
		log.Print("func: customerByCustomerNumberHandler")
	}

	url := strings.Split(r.URL.Path, "/")

	customerNumber, err := strconv.Atoi(url[5])
	if err != nil || customerNumber < 1 {

		response := &ResponseCustomerERROR{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "wrong_parameter",
		}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	if DebugConfiguration {
		log.Printf("Customer number: %d", customerNumber)
	}

	// Obtiene el cliente de la base de datos
	customer, err := GetCustomerByCustomerNumber(customerNumber)
	if err != nil {
		response := &ResponseCustomerERROR{
			Code:    http.StatusNotFound,
			Status:  "success",
			Message: "not_exist",
		}
		responseWriter(http.StatusNotFound, response, w)
		return
	}

	response := &ResponseCustomerOK{
		Code:    http.StatusOK,
		Status:  "success",
		Message: customer,
	}
	responseWriter(http.StatusOK, response, w)
}
