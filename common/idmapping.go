package common

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
				result.lastupdate = parsedValue
			}
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
IdMappingFromFilebase looks up a record by the filebase and returns it
*/
func IdMappingFromFilebase(config *Config, filebase string) (*IdMappingRecord, error) {
	ddbClient := config.GetDynamoClient()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	response, err := ddbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"filebase": &types.AttributeValueMemberS{Value: filebase}},
		TableName: &config.IdMappingTable,
	})
	if err != nil {
		return nil, err
	}

	if response.Item == nil {
		return nil, nil
	}

	return NewIdMappingRecord(&response.Item)
}

//func IdMappingFromOctid(config *Config, octid string) (*IdMappingRecord, error) {
//	ddbClient := config.GetDynamoClient()
//
//	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancelFunc()
//
//	//map[string]types.AttributeValue{"octopus_id": &types.AttributeValueMemberN{Value: octid}},
//	response, err := ddbClient.Query(ctx, &dynamodb.QueryInput{
//		IndexName:	aws.String("octid_index"),
//		TableName:                &config.IdMappingTable,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	if response.Item==nil {
//		return nil, nil
//	}
//
//	return NewIdMappingRecord(&response.Item)
//}
