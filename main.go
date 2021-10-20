package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

type Response struct {
	Health  string `json:"Health"`
	Message string `json:"Message"`
	Ok      bool   `json:"Ok"`
}

type Event struct {
	Name string `json:"name"`
}

func Handler(ctx context.Context, event Event) (Response, error) {
	message := "Failed"
	if event.Name == "getCSV" {
		message = getCSVData()
	}
	return Response{
		Health:  "UP",
		Message: fmt.Sprintf("Event Triggered %s", message),
		Ok:      true,
	}, nil
}

func getCSVData() string {
	return "Success"
}
