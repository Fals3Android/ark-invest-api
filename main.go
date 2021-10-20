package main

import (
	"fmt"
	"context"
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

func Handler(ctx context.Context, name Event) (Response, error) {
	return Response{
		Health:  "UP",
		Message: fmt.Sprintf("Event Triggered %s", name),
		Ok:      true,
	}, nil
}
