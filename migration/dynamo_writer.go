package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log"
	"reflect"
	"time"
	"strings"
)

type RawDynamoRecord map[string]types.AttributeValue

func buildPutRequests(queuePtr *[]*RawDynamoRecord) []types.WriteRequest {
	out := make([]types.WriteRequest, len(*queuePtr))
	for i, rec := range *queuePtr {
		out[i] = types.WriteRequest{PutRequest: &types.PutRequest{Item: *rec}}
	}
	return out
}

func commitQueue(ddbClient *dynamodb.Client, ctx context.Context, queuePtr *[]*RawDynamoRecord, tableNamePtr *string) error {
	var requestItems = map[string][]types.WriteRequest{
		*tableNamePtr: buildPutRequests(queuePtr),
	}

	for {
		log.Printf("INFO commitQueue committing %d records to dynamo", len(*queuePtr))
		req := &dynamodb.BatchWriteItemInput{
			RequestItems: requestItems,
		}

		response, err := ddbClient.BatchWriteItem(ctx, req)
		if err != nil {
			log.Printf("ERROR commitQueue could not commit %d records: %s", len(*queuePtr), err)
			return err
		}

		if len(response.UnprocessedItems) > 0 {
			requestItems = response.UnprocessedItems
		} else {
			log.Printf("INFO commitQueue completed")
			break
		}
	}
	return nil
}

func newWriteQueue() []*RawDynamoRecord {
	return make([]*RawDynamoRecord, 0)
}

/**
marshalGeneralRecord marshals a `GeneralRecord` structure from the sql reader into a DynamoDB record,
setting the types appropriately
*/
func marshalGeneralRecord(in *GeneralRecord, addUUID bool, nullableKeyFields string) (*RawDynamoRecord, error) {
	output := make(RawDynamoRecord, 0)
	for k, v := range *in {
		var val types.AttributeValue
		switch typeValue := v.(type) {
		case int8:
			val = &types.AttributeValueMemberBOOL{Value: typeValue != 0}
		case int16:
			val = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", typeValue)}
		case int32:
			val = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", typeValue)}
		case int64:
			val = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", typeValue)}
		case float64:
			val = &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", typeValue)}
		case string:
			if ((v == "") && (strings.Contains(nullableKeyFields, k))) {
				val = &types.AttributeValueMemberS{Value: "ABSENT"}
			} else {
				val = &types.AttributeValueMemberS{Value: typeValue}
			}
		case time.Time:
			val = &types.AttributeValueMemberS{Value: typeValue.Format(time.RFC3339)}
		default:
			return nil, errors.New(fmt.Sprintf("unknown type %s for value %v in column %s", reflect.TypeOf(v), v, k))
		}
		output[k] = val
	}
	if addUUID {
		uid, _ := uuid.NewRandom()
		output["uuid"] = &types.AttributeValueMemberS{Value: uid.String()}
	}
	return &output, nil
}

func AsyncDynamoWriter(inputCh chan GeneralRecord, ddbClient *dynamodb.Client, tableNamePtr *string, addUUID bool, nullableKeyFields string) chan error {
	errCh := make(chan error, 1)

	go func() {
		writeQueue := newWriteQueue()
		for {
			rec := <-inputCh
			if rec == nil {
				err := commitQueue(ddbClient, context.Background(), &writeQueue, tableNamePtr)
				errCh <- err //will be nil if there is no error
				return
			}

			ddbRec, err := marshalGeneralRecord(&rec, addUUID, nullableKeyFields)
			if err != nil {
				log.Printf("ERROR Could not marshal %v: %s", rec, err)
				errCh <- err
				return
			}

			writeQueue = append(writeQueue, ddbRec)
			if len(writeQueue) >= 25 {
				err := commitQueue(ddbClient, context.Background(), &writeQueue, tableNamePtr)
				if err != nil {
					errCh <- err
					return
				}
				writeQueue = newWriteQueue()
			}
		}
	}()
	return errCh
}
