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
	id          string
	date        string
	shares      int
	cusip       string
	market_value float64
	ticker      string
	fund        string
	weight      float64
	company     string
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
	
	entry := Entry{
		id:          id.String(),
		date:        "10/23/21",
		shares:      999,
		cusip:       "HI9409I8",
		market_value: 12351234.12,
		ticker:      "DWAC",
		fund:        "test",
		weight:      12.43,
		company:     "TEST",
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

	fmt.Println("Successfully added '" + entry.ticker + " 'to table': " + tableName)
}
