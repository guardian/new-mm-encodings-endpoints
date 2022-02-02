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
	"time"
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
	QueryEncodingsForFCSId(ctx context.Context, fcsid string) ([]*Encoding, error)
	QueryEncodingsForContentId(ctx context.Context, contentid int32, maybeSince *time.Time) ([]*Encoding, error)
	QueryIdMappings(ctx context.Context, indexName string, keyFieldName string, searchTerm interface{}) (*IdMappingRecord, error)
	GetAllMimeEquivalents(ctx context.Context) ([]*MimeEquivalent, error)
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

/*
Internal function that takes a QueryOutput and builds a list of Encodings to return then sorts them by VBitrate
*/
func _marshalResponseToSortedEncodings(response *dynamodb.QueryOutput) ([]*Encoding, error) {
	var err error
	encodings := make([]*Encoding, len(response.Items))
	for i, rawData := range response.Items {
		encodings[i], err = EncodingFromDynamo((*RawDynamoRecord)(&rawData))
		if err != nil {
			log.Printf("ERROR QueryEncodingsForFCSId could not marshal item %d (%v): %s", i, rawData, err)
			return nil, err
		}
	}

	sort.Slice(encodings, func(i int, j int) bool {
		return encodings[i].VBitrate > encodings[j].VBitrate
	})
	return encodings, nil
}

/*
QueryEncodingsForFCSId searches the Encodings table for videos corresponding to the given fcsid
and returns a slice of pointers to the marshalled Encoding objects
*/
func (ops *DynamoDbOpsImpl) QueryEncodingsForFCSId(ctx context.Context, fcsid string) ([]*Encoding, error) {
	//equivalent SQL is select * from encodings left join mime_equivalents on (real_name=encodings.format) where fcs_id='$fcsid' order by vbitrate desc
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("fcs_id").Equal(expression.Value(fcsid))).
		Build()
	if err != nil {
		log.Printf("ERROR QueryEncodingsForFCSId could not build the query expression: %s", err)
		return nil, err
	}

	rq := &dynamodb.QueryInput{
		TableName:                 ops.config.EncodingsTablePtr(),
		ExclusiveStartKey:         nil,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	response, err := ops.client.Query(ctx, rq)
	if err != nil {
		log.Printf("ERROR QueryEncodingsForFCSId could not perform the query: %s", err)
		return nil, err
	}

	return _marshalResponseToSortedEncodings(response)
}

/*
QueryEncodingsForContentId looks up records from the Encodings table corresponding to the given contentid.
This is the "fallback mode" query for when no title version id can be found

Arguments:
- ctx - context that can be used to cancel this operation
- contentid - the content ID to query for
- maybeSince - nullable pointer to a time. If this is non-NULL then only records with 'lastupdate' equal to or since this time will be retrieved
Returns:
- a slice of pointers to matching Encoding records on success
- an error on failure
*/
func (ops *DynamoDbOpsImpl) QueryEncodingsForContentId(ctx context.Context, contentid int32, maybeSince *time.Time) ([]*Encoding, error) {
	//equivalent SQL is select * from encodings left join mime_equivalents on (real_name=encodings.format) where contentid=$contentid order by vbitrate desc,lastupdate desc
	keyTerms := expression.Key("contentid").Equal(expression.Value(contentid))
	if maybeSince != nil {
		keyTerms = keyTerms.And(expression.Key("lastupdate").GreaterThanEqual(expression.Value(maybeSince.Format(time.RFC3339))))
	}

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyTerms).
		Build()

	if err != nil {
		log.Printf("ERROR QueryEncodingsForContentId could not build the query expression: %s", err)
		return nil, err
	}
	rq := &dynamodb.QueryInput{
		TableName:                 ops.config.EncodingsTablePtr(),
		ExclusiveStartKey:         nil,
		IndexName:                 aws.String("contentid"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	response, err := ops.client.Query(ctx, rq)
	if err != nil {
		log.Printf("ERROR QueryEncodingsForContentId could not execute the query: %s", err)
		return nil, err
	}

	encodings, err := _marshalResponseToSortedEncodings(response)
	if err == nil {
		//apply a most-recent-first search
		sort.Slice(encodings, func(i int, j int) bool {
			return encodings[i].LastUpdate.Unix() > encodings[j].LastUpdate.Unix()
		})
		return encodings, nil
	} else {
		return nil, err
	}
}

const IdMappingIndexFilebase = "filebase"
const IdMappingKeyfieldFilebase = "filebase"
const IdMappingIndexOctid = "octopusid"
const IdMappingKeyfieldOctid = "octopus_id"

/*
QueryIdMappings performs a lookup on the IdMappings table.  There should only ever be 1 or 0 matches; in the event of
more than one a warning is logged and the first value is used.

Arguments:
- ctx - context that can be used to cancel the operation, normally passed through from lambda
- indexName - the index to query. Should be an IdMappingIndex* const that is defined in `common`
- keyFieldName - the key field to query. Should be the IdMappingKeyfield* const that corresponds to the given IdMappingIndex* used for indexName
- searchTerm - the value to search. This must be compatible with the field type or Dynamo will return a runtime error
Returns:
- a pointer to an IdMappingRecord on success, or nil if nothing found or error. If both return values are `nil` that means that
there was no data found
- an error on failure
*/
func (ops *DynamoDbOpsImpl) QueryIdMappings(ctx context.Context, indexName string, keyFieldName string, searchTerm interface{}) (*IdMappingRecord, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.
			Key(keyFieldName).
			Equal(expression.Value(searchTerm))).
		Build()

	if err != nil {
		return nil, err
	}

	response, err := ops.client.Query(ctx, &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		TableName:                 aws.String(ops.config.IdMappingTable()),
		IndexName:                 &indexName,
		Limit:                     aws.Int32(20),
	})

	if err != nil {
		return nil, err
	}

	if response.Items == nil {
		return nil, nil
	}

	if len(response.Items) == 0 {
		return nil, nil
	} else if len(response.Items) == 1 {
		return NewIdMappingRecord(&response.Items[0])
	} else {
		log.Printf("WARNING Got %d idmapping records, expected only 1. Note that there is a hard limit of 20. Using the most recent.", len(response.Items))
		mostRecent := len(response.Items) - 1 //the indexing setup in DynamoDB returns the oldest first (sort key is lastupdate)
		return NewIdMappingRecord(&response.Items[mostRecent])
	}
}

/*
GetAllMimeEquivalents downloads the MIME equivalents table to use for lookups
*/
func (ops *DynamoDbOpsImpl) GetAllMimeEquivalents(ctx context.Context) ([]*MimeEquivalent, error) {
	rq := &dynamodb.ScanInput{
		TableName: ops.config.MimeEquivalentsTablePtr(),
	}
	response, err := ops.client.Scan(ctx, rq)
	if err != nil {
		log.Printf("ERROR Can't load in mime equivalents: %s", err)
		return nil, err
	}

	results := make([]*MimeEquivalent, len(response.Items))
	for i, raw := range response.Items {
		results[i], err = MimeEquivalentFromDynamo((*RawDynamoRecord)(&raw))
		if err != nil {
			log.Printf("ERROR Can't load in record %d from mime equivalents (%v): %s", i, raw, err)
			return nil, err
		}
	}
	return results, nil
}
