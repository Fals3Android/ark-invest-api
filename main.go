package main

import (
	"context"
	"fmt"
	"encoding/json"
	"log"
	"os"
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
	logEventData(event)

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

func logEventData(event Event) {
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("EVENT: %s", eventJson)
	log.Printf("REGION: %s", os.Getenv("AWS_REGION"))
	log.Println("ALL ENV VARS:")
	for _, element := range os.Environ() {
		log.Println(element)
	}
}
