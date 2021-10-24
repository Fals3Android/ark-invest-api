package main

import (
	"context"
	"fmt"
	"encoding/json"
	"log"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"encoding/csv"
	"io"
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
		csv, status := readCSVFromUrl("https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_INNOVATION_ETF_ARKQ_HOLDINGS.csv")
		message = status
		log.Printf("DATA: %s", csv)
	}
	
	return Response{
		Health:  "UP",
		Message: fmt.Sprintf("Event Triggered %s", message),
		Ok:      true,
	}, nil
}

func readCSVFromUrl(url string) ([][]string, string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Chrome/56.0.2924.76")
	response, err := client.Do(req)
	if err != nil {
		return [][]string{}, "getCSV() Failed"
	}

	defer response.Body.Close()
	reader := csv.NewReader(response.Body)
	list := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
			// should or do something useful log.Fatal(err)
		}
		list = append(list, record)
	}

	return list, "getCSV() Success"
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
