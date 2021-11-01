package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Entry struct {
	Id          string  `json:"id"`
	Date        string  `json:"date"`
	Fund        string  `json:"fund"`
	Company     string  `json:"company"`
	Ticker      string  `json:"ticker"`
	Cusip       string  `json:"cusip"`
	Shares      int     `json:"shares"`
	MarketValue float64 `json:"market_value"`
	Weight      float64 `json:"weight"`
}

func putBatchRequest(list [][]string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	tableName := "ARK_INNOVATION_ETF_ARKQ_HOLDINGS"
	svc := dynamodb.New(sess)
	entries := convertRowsToAttributes(list)
	batchRequestItems := getBatchRequestItems(entries, tableName)
	fmt.Println(batchRequestItems)

	for _, item := range batchRequestItems {
		input := &dynamodb.BatchWriteItemInput{
			RequestItems: item,
		}
		result, err := svc.BatchWriteItem(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeProvisionedThroughputExceededException:
					fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				case dynamodb.ErrCodeResourceNotFoundException:
					fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
					fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				case dynamodb.ErrCodeRequestLimitExceeded:
					fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println(result, "hello")
	}
}

func getBatchRequestItems(list []map[string]*dynamodb.AttributeValue, tableName string) []map[string][]*dynamodb.WriteRequest {
	batchedRequests := make([]map[string][]*dynamodb.WriteRequest, 0) // you can send endless batches
	request := make(map[string][]*dynamodb.WriteRequest)              // each request must not exceed 25 writes
	writeRequests := []*dynamodb.WriteRequest{}

	for _, item := range list {
		if len(writeRequests) == 25 {
			request[tableName] = writeRequests
			batchedRequests = append(batchedRequests, request)
			request = make(map[string][]*dynamodb.WriteRequest)
			writeRequests = []*dynamodb.WriteRequest{}
		}

		// row, err := dynamodbattribute.MarshalMap(item)
		// if err != nil {
		// 	log.Fatalf("Got error calling PutItem: %s", err)
		// }

		writeRequests = append(writeRequests, &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{
			Item: item,
		}})
	}

	if len(writeRequests) > 0 {
		request[tableName] = writeRequests
		batchedRequests = append(batchedRequests, request)
	}

	return batchedRequests
}

func convertRowsToAttributes(list [][]string) []map[string]*dynamodb.AttributeValue {
	attributeList := make([]map[string]*dynamodb.AttributeValue, 0)
	for i := 1; i < len(list); i++ {
		item := list[i]
		id := uuid.New()
		shares := strings.ReplaceAll(item[5], ",", "")
		marketValue := strings.ReplaceAll(strings.TrimPrefix(item[6], "$"), ",", "")
		weight := strings.ReplaceAll(strings.TrimSuffix(item[7], "%"), ",", "")
		row := map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id.String()),
			},
			"Date": {
				S: aws.String(item[0]),
			},
			"Fund": {
				S: aws.String(item[1]),
			},
			"Company": {
				S: aws.String(item[2]),
			},
			"Ticker": {
				S: aws.String(item[3]),
			},
			"Cusip": {
				S: aws.String(item[4]),
			},
			"Shares": {
				N: aws.String(shares),
			},
			"MarketValue": {
				N: aws.String(marketValue),
			},
			"Weight": {
				N: aws.String(weight),
			},
		}
		attributeList = append(attributeList, row)
	}
	return attributeList
}

func convertRowsToEntries(list [][]string) []Entry {
	entries := []Entry{}
	for i := 1; i < len(list); i++ {
		item := list[i]
		id := uuid.New()
		shares, _ := strconv.Atoi(item[4])
		marketValue, _ := strconv.ParseFloat(item[6], 64)
		weight, _ := strconv.ParseFloat(item[7], 64)

		entries = append(entries, Entry{
			Id:          id.String(),
			Date:        item[0],
			Fund:        item[1],
			Company:     item[2],
			Ticker:      item[3],
			Shares:      shares,
			Cusip:       item[5],
			MarketValue: marketValue,
			Weight:      weight,
		})
	}
	return entries
}

func createSingleDBItemExample(list [][]string) {
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
