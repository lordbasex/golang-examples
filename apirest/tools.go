package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func responseWriter(status int, data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func stringToBool(str string) (bool, error) {
	lowerStr := strings.ToLower(str)

	if lowerStr == "true" {
		return true, nil
	} else if lowerStr == "false" {
		return false, nil
	} else {
		return false, fmt.Errorf("Invalid boolean string")
	}
}
