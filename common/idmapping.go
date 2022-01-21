package common

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"strconv"
	"time"
)

type IdMappingRecord struct {
	contentId  string
	filebase   string //base index
	project    *string
	lastupdate time.Time //range key for all indices
	octopus_id *int64    //indexed
}

func NewIdMappingRecord(from *map[string]types.AttributeValue) (*IdMappingRecord, error) {
	var result IdMappingRecord
	if contentId, haveContentId := (*from)["contentId"]; haveContentId {
		if contentIdString, contentIdIsString := contentId.(*types.AttributeValueMemberS); contentIdIsString {
			result.contentId = contentIdString.Value
		}
	}
	if filebase, haveFileBase := (*from)["filebase"]; haveFileBase {
		if filebaseString, filebaseIsString := filebase.(*types.AttributeValueMemberS); filebaseIsString {
			result.filebase = filebaseString.Value
		}
	}
	if project, haveProject := (*from)["project"]; haveProject {
		if projectString, projectIsString := project.(*types.AttributeValueMemberS); projectIsString {
			copiedValue := projectString.Value
			result.project = &copiedValue
		}
	}
	if lastupdate, haveLastUpdate := (*from)["lastupdate"]; haveLastUpdate {
		if lastUpdateString, lastUpdateIsString := lastupdate.(*types.AttributeValueMemberS); lastUpdateIsString {
			parsedValue, err := time.Parse(time.RFC3339, lastUpdateString.Value)
			if err != nil {
				return nil, err
			}
			result.lastupdate = parsedValue
		}
	}
	if octId, haveOctid := (*from)["octopus_id"]; haveOctid {
		if octIdNum, octIdIsNum := octId.(*types.AttributeValueMemberN); octIdIsNum {
			intValue, _ := strconv.ParseInt(octIdNum.Value, 10, 64)
			result.octopus_id = &intValue
		}
	}

	nullTime := time.Time{}
	if result.contentId == "" || result.filebase == "" || result.lastupdate == nullTime {
		return nil, errors.New("ID mapping record is inaccurate, does not contain required fields")
	}
	return &result, nil
}

/**
internalDbLookup is an internal method that performs an 'equals' query against an index on a dynemo table

Parameters:
ddbClient - dynamodb client instance
tableName - string representing the table to query
indexName - string representing the index to query
keyFieldName - name of the hash key on the index to query
keyValue  - the value to match against
limit     - maximum number of rows to return
*/
func internalDbLookup(ctx context.Context, ddbClient dynamodb.QueryAPIClient, tableName string, indexName string, keyFieldName string, keyValue interface{}, limit int) (*dynamodb.QueryOutput, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.
			Key(keyFieldName).
			Equal(expression.Value(keyValue))).
		Build()

	if err != nil {
		return nil, err
	}

	return ddbClient.Query(ctx, &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		TableName:                 &tableName,
		IndexName:                 &indexName,
		Limit:                     aws.Int32(int32(limit)),
	})
}

/**
formatResponse takes the QueryOutput from dynamo and marshals the results into the first IdMappingRecord, if present.
Emits a warning if there is more than one record in the response.
If not, then returns nil.
On error, returns the error
*/
func formatResponse(response *dynamodb.QueryOutput, err error) (*IdMappingRecord, error) {
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
		log.Printf("WARNING Got %d idmapping records, expected only 1. Note that there is a hard limit of 20.", len(response.Items))
		return NewIdMappingRecord(&response.Items[0])
	}
}

/*
IdMappingFromFilebase looks up a record by the filebase and returns it
*/
func IdMappingFromFilebase(ctx context.Context, config Config, filebase string) (*IdMappingRecord, error) {
	ddbClient := config.GetDynamoClient()

	response, err := internalDbLookup(ctx, ddbClient, config.IdMappingTable(), "filebase", "filebase", filebase, 20)

	return formatResponse(response, err)
}

/*
IdMappingFromOctid looks up a record by the octopus id and returns it
*/
func IdMappingFromOctid(ctx context.Context, config Config, octid int64) (*IdMappingRecord, error) {
	ddbClient := config.GetDynamoClient()

	response, err := internalDbLookup(ctx, ddbClient, config.IdMappingTable(), "octopusid", "octopus_id", octid, 20)

	return formatResponse(response, err)
}
