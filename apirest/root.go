package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// root
func root(w http.ResponseWriter, r *http.Request) {
	if DebugConfiguration {
		log.Print("func: root")
	}

	version, _ := GetDataBase()

	location_currentzone, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	currentTime := time.Now()
	timestamp := currentTime.Unix()

	response := &Root{
		Project:   "api-ivr-autogestion",
		Status:    http.StatusOK,
		GoLang:    "1.20.2",
		DataBase:  version,
		Version:   "1.1",
		Timezone:  fmt.Sprintf("%v", location_currentzone),
		Timestamp: timestamp,
	}
	responseWriter(http.StatusOK, response, w)
}
