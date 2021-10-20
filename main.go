package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

type Request struct {
	ID    int    `json:"ID"`
	Value string `json:"Value"`
}

type Response struct {
	Health  string `json:"Health"`
	Message string `json:"Message"`
	Ok      bool   `json:"Ok"`
}

type Event struct {
	Name string `json:"name"`
}

func Handler(request Request, name Event) (Response, error) {
	return Response{
		Health:  "UP",
		Message: fmt.Sprintf("Event Triggered %s", name.Name),
		Ok:      true,
	}, nil
}
