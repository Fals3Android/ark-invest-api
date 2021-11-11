package main

import (
	"context"
	"fmt"
	"log"

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

	message := "No Event Passed"

	if event.Name == "getCSV" {
		csv, status := getCSVFromUrl("https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_AUTONOMOUS_TECH._&_ROBOTICS_ETF_ARKQ_HOLDINGS.csv")
		message = status
		log.Printf("DATA: %s", csv)
		log.Printf("Writing To DynamoDB")
		tableName := "ARK_AUTONOMOUS_TECH_AND_ROBOTICS_ETF_ARKQ_HOLDINGS"
		putBatchRequest(csv, tableName)
	}

	return Response{
		Health:  "UP",
		Message: fmt.Sprintf("Event Triggered %s", message),
		Ok:      true,
	}, nil
}
