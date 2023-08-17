package main

import (
	"fmt"
	"log"
	"net/http"
)

// helloWorld
func helloWorld(w http.ResponseWriter, r *http.Request) {
	if DebugConfiguration {
		log.Print("func: helloWorld")
	}

	response := &Msg{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Hello World " + fmt.Sprintf("%v", UserValues[1]),
	}
	responseWriter(http.StatusOK, response, w)
}
