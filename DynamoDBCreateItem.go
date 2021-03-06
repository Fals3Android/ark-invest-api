package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

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

func putBatchRequest(funds map[string]MetaInfo) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	batchedRequests := make([]map[string][]*dynamodb.WriteRequest, 0)
	for _, value := range funds {
		batchedRequests = append(batchedRequests, getBatchRequestItems(value.attributeValues, value.tableName)...)
	}

	// batchGroups := consolidateBatchGroups(batchedRequests)
	var wg sync.WaitGroup

	for _, item := range batchedRequests {
		wg.Add(1)
		go doBatchWrite(&wg, sess, item)
	}

	wg.Wait() //Works ! writes new values to the db but for some reason causes timeouts
}

func doBatchWrite(wg *sync.WaitGroup, sess *session.Session, item map[string][]*dynamodb.WriteRequest) {
	defer wg.Done()
	svc := dynamodb.New(sess)

	input := &dynamodb.BatchWriteItemInput{
		RequestItems:                item,
		ReturnConsumedCapacity:      aws.String("INDEXES"),
		ReturnItemCollectionMetrics: aws.String("SIZE"),
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
	fmt.Println(result)
}

func consolidateBatchGroups(batchedRequests []map[string][]*dynamodb.WriteRequest) []map[string][]*dynamodb.WriteRequest {
	groups := make([]map[string][]*dynamodb.WriteRequest, 0)
	uniq := make(map[string][]*dynamodb.WriteRequest)
	counter := 0
	for len(batchedRequests) != 0 {
		fmt.Println(batchedRequests[counter], counter, uniq, batchedRequests)
		current := batchedRequests[counter]

		for key, value := range current {
			_, ok := uniq[key]
			if !ok {
				uniq[key] = value
				batchedRequests = append(batchedRequests[:counter], batchedRequests[counter+1:]...)
				counter = 0
				continue
			}
		}

		if counter >= len(batchedRequests)-1 {
			groups = append(groups, uniq)
			uniq = make(map[string][]*dynamodb.WriteRequest)
			counter = 0
		} else {
			counter++
		}
	}
	return groups
}

func getBatchRequestItems(list []map[string]*dynamodb.AttributeValue, tableName string) []map[string][]*dynamodb.WriteRequest {
	batch := make([]map[string][]*dynamodb.WriteRequest, 0)
	request := make(map[string][]*dynamodb.WriteRequest) // each request must not exceed 25 writes
	writeRequests := []*dynamodb.WriteRequest{}

	for _, item := range list {
		if len(writeRequests) == 24 {
			request[tableName] = writeRequests
			batch = append(batch, request)
			request = make(map[string][]*dynamodb.WriteRequest)
			writeRequests = []*dynamodb.WriteRequest{}
		}

		writeRequests = append(writeRequests, &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{
			Item: item,
		}})
	}

	if len(writeRequests) > 0 {
		request[tableName] = writeRequests
		batch = append(batch, request)
	}
	return batch
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
			"id": {
				S: aws.String(id.String()),
			},
			"date": {
				S: aws.String(item[0]),
			},
			"fund": {
				S: aws.String(item[1]),
			},
			"company": {
				S: aws.String(item[2]),
			},
			"ticker": {
				S: aws.String(item[3]),
			},
			"cusip": {
				S: aws.String(item[4]),
			},
			"shares": {
				N: aws.String(shares),
			},
			"market_value": {
				N: aws.String(marketValue),
			},
			"weight": {
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
