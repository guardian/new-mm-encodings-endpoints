package common

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"strconv"
	"time"
)

type IdMappingRecord struct {
	contentId  int32
	filebase   string //base index
	project    *string
	lastupdate time.Time //range key for all indices
	octopus_id *int64    //indexed
}

func NewIdMappingRecord(from *map[string]types.AttributeValue) (*IdMappingRecord, error) {
	log.Printf("DEBUG NewIdMappingRecord from %v", from)
	var result IdMappingRecord
	if contentId, haveContentId := (*from)["contentid"]; haveContentId {
		if contentIdNumber, contentIdIsNum := contentId.(*types.AttributeValueMemberN); contentIdIsNum {
			intValue, err := strconv.ParseInt(contentIdNumber.Value, 10, 32)
			if err != nil {
				return nil, err
			}
			result.contentId = int32(intValue)
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
	if result.contentId == 0 || result.filebase == "" || result.lastupdate == nullTime {
		log.Printf("ERROR ID mapping record is inaccurate, does not contain required fields")
		log.Printf("ERROR Partial result was %d %s %s", result.contentId, result.filebase, result.lastupdate)
		return nil, errors.New("ID mapping record is inaccurate, does not contain required fields")
	}
	return &result, nil
}
