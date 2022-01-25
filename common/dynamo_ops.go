package common

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
	"sort"
)

type DynamoDbOpsImpl struct {
	client *dynamodb.Client
	config Config
}

/*
DynamoDbOps abstracts the actual DynamoDb operations so that we can mock them in testing
*/
type DynamoDbOps interface {
	QueryFCSIdForContentId(ctx context.Context, contentId int32) (*[]string, error)
}

/*
NewDynamoDbOps creates a new DynamoDbOps object from the given configuration
*/
func NewDynamoDbOps(config Config) DynamoDbOps {
	return &DynamoDbOpsImpl{client: config.GetDynamoClient(), config: config}
}

type SortableString struct {
	StringValue string
	LastUpdate  string
}

/*
QueryFCSIdForContentId queries the Encodings table for all FCS IDs relating to the given `contentId`
Arguments:
- ctx - context that can be used to cancel the operation
- contentId - contentId to query
Returns:
- a pointer to a list of FCS ID values, or null on error
- an error if the operation fails, or null on success
This does the equivalent of
```sql
select fcs_id from encodings left join mime_equivalents on (real_name=encodings.format)where contentid=$contentid order by lastupdate desc
```
from line 273 of the original code
*/
func (ops *DynamoDbOpsImpl) QueryFCSIdForContentId(ctx context.Context, contentId int32) (*[]string, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("contentid").Equal(expression.Value(contentId))).
		Build()

	if err != nil {
		log.Printf("ERROR Could not build query expression for FCS ID -> Content ID: %s", err)
		return nil, err
	}

	output := make([]SortableString, 0)
	var nextStartKey map[string]types.AttributeValue
	ctr := 0
	for {
		//query the contentid index to get the FCS IDs
		rq := &dynamodb.QueryInput{
			TableName:                 ops.config.EncodingsTablePtr(),
			ExclusiveStartKey:         nextStartKey,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			IndexName:                 aws.String("contentid"),
			KeyConditionExpression:    expr.KeyCondition(),
		}
		results, err := ops.client.Query(ctx, rq)
		if err != nil {
			log.Printf("ERROR FCS ID -> Content ID query failed on iteration %d: %s", ctr, err)
			return nil, err
		}

		for _, result := range results.Items {
			output = append(output, SortableString{
				StringValue: extractDynamoField((*RawDynamoRecord)(&result), "fcs_id", reflect.String, true).(string),
				LastUpdate:  extractDynamoField((*RawDynamoRecord)(&result), "lastupdate", reflect.String, true).(string),
			})
		}

		ctr++
		nextStartKey = results.LastEvaluatedKey
		if nextStartKey == nil {
			break
		}
	}

	//rely on lexographical properties of the iso timestamp to do the date sort
	//this does a most-recent-first sort
	sort.Slice(output, func(i int, j int) bool {
		return output[j].LastUpdate < output[i].LastUpdate
	})

	finalOutputs := make([]string, len(output))
	for i, v := range output {
		log.Printf("DEBUG FCS query for contentid %d got %v @ %v", contentId, v.StringValue, v.LastUpdate)
		if v.StringValue != "" {
			finalOutputs[i] = v.StringValue
		}
	}
	return &finalOutputs, nil
}
