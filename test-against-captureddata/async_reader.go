package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/guardian/new-encodings-endpoints/common"
	"log"
)

func AsyncRecordReader(client *dynamodb.Client, tableName *string, pageSize int32) (chan *EndpointEvent, chan error) {
	outputCh := make(chan *EndpointEvent, pageSize*2)
	errCh := make(chan error, 1)

	go func() {
		var continuationKey map[string]types.AttributeValue
		for {
			req := &dynamodb.ScanInput{
				TableName:         tableName,
				ExclusiveStartKey: continuationKey,
				Limit:             aws.Int32(pageSize),
			}
			log.Printf("Retrieving from %s...", *tableName)
			response, err := client.Scan(context.Background(), req)
			if err != nil {
				log.Printf("ERROR %s", err)
				errCh <- err
				close(outputCh)
				return
			}
			log.Printf("DEBUG Got page of %d items from %s", len(response.Items), *tableName)
			for _, item := range response.Items {
				event, marshalErr := EndpointEventFromDynamo((*common.RawDynamoRecord)(&item))
				if marshalErr != nil {
					log.Printf("ERROR %s", marshalErr)
					errCh <- err
					close(outputCh)
					return
				}
				outputCh <- event
			}
			if response.LastEvaluatedKey == nil {
				break //docs say that LastEvaluatedKey is blank when we get to the end
			} else {
				continuationKey = response.LastEvaluatedKey
			}
		}
		log.Printf("INFO AsyncReader reached the end of records")
		close(outputCh)
	}()
	return outputCh, errCh
}
