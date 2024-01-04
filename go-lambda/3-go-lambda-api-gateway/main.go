package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type myJSON struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func main() {
	lambda.Start(handler)
}

func handler(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var requestBody myJSON

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, DELETE, PUT, OPTIONS",
		},
	}

	err := json.Unmarshal([]byte(ev.Body), &requestBody)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	response.Body = fmt.Sprintf(`{"message:" "Hello %s %s"}`, requestBody.Name, requestBody.Lastname)
	response.StatusCode = http.StatusOK

	return response, nil

}
