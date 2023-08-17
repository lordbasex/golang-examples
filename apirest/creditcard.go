package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func CreateCreditCardHandler(w http.ResponseWriter, r *http.Request) {

	body := &CreateCredit{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		response := &Msg{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "parameters_are_required_in_the_body",
		}
		responseWriter(http.StatusBadRequest, response, w)
		return

	}

	if validErrs := body.validateCreateCreditCard(); len(validErrs) > 0 {
		err := map[string]interface{}{
			"error":    false,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": validErrs}
		responseWriter(http.StatusBadRequest, err, w)
		return
	}

	// Call SetCreditCard function
	err := body.SetCreditCard()

	if err != nil {
		err := map[string]interface{}{
			"error":    true,
			"code":     http.StatusInternalServerError,
			"status":   "error",
			"response": "No_created_payment"}
		responseWriter(http.StatusBadRequest, err, w)
		return
	}

	okResponse := map[string]interface{}{
		"error":    false,
		"code":     http.StatusOK,
		"status":   "success",
		"response": "Payment_created_successfully",
	}
	responseWriter(http.StatusOK, okResponse, w)
}

// Function to check if a string is numeric
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func SetIntToString(i int) string {
	r := strconv.Itoa(i)
	return r
}

func CreateCreditAll(w http.ResponseWriter, r *http.Request) {

	search, err := GetCreditCardAll()

	if err != nil {
		err := map[string]interface{}{
			"error":   true,
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "No_created_payment"}
		responseWriter(http.StatusBadRequest, err, w)
		return
	}

	okResponse := map[string]interface{}{
		"error":   false,
		"code":    http.StatusOK,
		"status":  "success",
		"message": search,
	}
	responseWriter(http.StatusOK, okResponse, w)
}

func CreateCreditSearch(w http.ResponseWriter, r *http.Request) {

	url := strings.Split(r.URL.Path, "/")

	uniqueIDStr := url[5] // Extract the part of the URL path you need

	if uniqueIDStr == "" {
		response := &ResponseCustomerERROR{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "wrong_parameter",
		}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	search, err := GetCreditCardByUniqueID(uniqueIDStr)

	if err != nil {
		err := map[string]interface{}{
			"error":   true,
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "No_created_payment"}
		responseWriter(http.StatusBadRequest, err, w)
		return
	}

	okResponse := map[string]interface{}{
		"error":   false,
		"code":    http.StatusOK,
		"status":  "success",
		"message": search,
	}
	responseWriter(http.StatusOK, okResponse, w)
}
