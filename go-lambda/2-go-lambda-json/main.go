package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type myJSON struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func main() {
	lambda.Start(handler)
}

func handler(e myJSON) (string, error) {

	log.Println("Hello Lambda :)")

	response := fmt.Sprintf("Hello: %s %s", e.Name, e.Lastname)

	return response, nil
}
