package main

import (
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
	Message string `json:"Message"`
	Ok      bool   `json:"Ok"`
}

func Handler(request Request) (Response, error) {
	return Response{
		Message: "UP",
		Ok:      true,
	}, nil
}
