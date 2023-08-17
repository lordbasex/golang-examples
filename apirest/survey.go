package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func CreateSurvey(w http.ResponseWriter, r *http.Request) {

	body := &Survey{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "parameters_are_required_in_the_body"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	if validErrs := body.validateCreateSurvey(); len(validErrs) > 0 {
		err := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": validErrs}
		responseWriter(http.StatusBadRequest, err, w)
		return
	}

	// Call SetCreateSurvey function
	err := body.SetCreateSurvey()

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

func UpdateSurvey(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")

	uniqueIDStr := url[5] // Extract the part of the URL path you need

	if uniqueIDStr == "" {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "wrong_parameter"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	var update SurveyUpdate
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&update); err != nil {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "invalid_json"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	// Validate that only one of q1, q2, q3 is present
	count := 0
	if update.Q1 != 0 {
		count++
	}
	if update.Q2 != 0 {
		count++
	}
	if update.Q3 != 0 {
		count++
	}

	if count != 1 {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "invalid_question_count"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	var chosenQuestion string
	var chosenValue int

	if update.Q1 != 0 {
		chosenQuestion = "q1"
		chosenValue = update.Q1
	} else if update.Q2 != 0 {
		chosenQuestion = "q2"
		chosenValue = update.Q2
	} else if update.Q3 != 0 {
		chosenQuestion = "q3"
		chosenValue = update.Q3
	} else {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "no_question_provided"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	// Validate chosenValue is a single digit (0 to 9)
	if chosenValue < 0 || chosenValue > 9 {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "invalid_chosen_value"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	if DebugConfiguration {
		log.Printf("uniqueID: %s question: %s Value: %d\n", uniqueIDStr, chosenQuestion, chosenValue)
	}

	err := SetUpdateSurvey(uniqueIDStr, chosenQuestion, chosenValue)
	if err != nil {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "invalid_json"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	} else {
		okResponse := map[string]interface{}{
			"error":    false,
			"code":     http.StatusOK,
			"status":   "success",
			"response": "survey_update_successfully",
		}
		responseWriter(http.StatusOK, okResponse, w)
	}
}

func AbortedSurvey(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")

	uniqueIDStr := url[5] // Extract the part of the URL path you need

	if uniqueIDStr == "" {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "wrong_parameter"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	}

	err := SetAbortedSurvey(uniqueIDStr)
	if err != nil {
		response := map[string]interface{}{
			"error":    true,
			"code":     http.StatusBadRequest,
			"status":   "error",
			"response": "invalid_json"}
		responseWriter(http.StatusBadRequest, response, w)
		return
	} else {
		okResponse := map[string]interface{}{
			"error":    false,
			"code":     http.StatusOK,
			"status":   "success",
			"response": "survey_aboort_successfully",
		}
		responseWriter(http.StatusOK, okResponse, w)
	}

}
