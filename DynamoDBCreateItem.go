package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Entry struct {
	Id          string
	Date        string
	Shares      int
	Cusip       string
	MarketValue float64
	Ticker      string
	Fund        string
	Weight      float64
	Company     string
}

func createDBItem(list [][]string) {
	// cfg := aws.Config{
	// 	Endpoint:   aws.String("http://localhost:8000"),
	// 	Region:     aws.String("us-west-2"),
	// 	MaxRetries: aws.Int(12),
	// }
	// sess := session.Must(session.NewSession())

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	
	svc := dynamodb.New(sess)
	id := uuid.New()
	fmt.Println(id.String())
	entry := Entry{
		Id:          id.String(),
		Date:        "10/23/21",
		Shares:      999,
		Cusip:       "HI9409I8",
		MarketValue: 12351234.12,
		Ticker:      "DWAC",
		Fund:        "test",
		Weight:      12.43,
		Company:     "TEST",
	}

	av, err := dynamodbattribute.MarshalMap(entry)
	if err != nil {
		log.Fatalf("Got error marshalling new item: %s", err)
	}

	tableName := "ARK_INNOVATION_ETF_ARKQ_HOLDINGS"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added '" + entry.Ticker + " 'to table': " + tableName)
}
