package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	log.Print("Start Lambda")
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Hello Lambda :)")

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
	}

	return response, nil
}
