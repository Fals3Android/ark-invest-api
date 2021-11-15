package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestConsolidateBatchGroups(t *testing.T) {
	stub := []map[string][]*dynamodb.WriteRequest{
		{
			"tableOne": []*dynamodb.WriteRequest{},
		},
		{
			"tableOne": []*dynamodb.WriteRequest{},
		},
		{
			"tableTwo": []*dynamodb.WriteRequest{},
		},
		{
			"tableThree": []*dynamodb.WriteRequest{},
		},
		{
			"tableThree": []*dynamodb.WriteRequest{},
		},
		{
			"tableFour": []*dynamodb.WriteRequest{},
		},
	}
	mock := []map[string][]*dynamodb.WriteRequest{
		{
			"tableOne":   []*dynamodb.WriteRequest{},
			"tableTwo":   []*dynamodb.WriteRequest{},
			"tableThree": []*dynamodb.WriteRequest{},
			"tableFour":  []*dynamodb.WriteRequest{},
		},
		{
			"tableOne":   []*dynamodb.WriteRequest{},
			"tableThree": []*dynamodb.WriteRequest{},
		},
	}
	result := consolidateBatchGroups(stub)
	if len(result) != 2 {
		t.Error("The result length does not match")
	}
	for index, elem := range result {
		for key, _ := range elem {
			_, ok := mock[index][key]
			if !ok {
				t.Error("Table consolidation has failed")
			}
		}
	}
}
