package main

import (
	"context"
	"fmt"
	"sync"

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

func handleDataTransaction(fund MetaInfo) [][]string {
	csv, _ := getCSVFromUrl(fund.url)
	return csv
}

func Handler(ctx context.Context, event Event) (Response, error) {
	logEventData(event)

	if event.Name == "getCSV" {
		var wg sync.WaitGroup
		funds := getTablesAndUrls()

		for key, value := range funds {
			ch := make(chan [][]string)
			wg.Add(1)
			go func(value MetaInfo) {
				defer wg.Done()
				ch <- handleDataTransaction(value)
			}(value)
			csv := <-ch
			value.data = csv
			value.attributeValues = convertRowsToAttributes(csv)
			funds[key] = value
		}

		wg.Wait()
		putBatchRequest(funds)
	}

	return Response{
		Health:  "UP",
		Message: fmt.Sprintf("Event Triggered %s", event.Name),
		Ok:      true,
	}, nil
}
