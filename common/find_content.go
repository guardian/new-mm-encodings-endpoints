package common

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"regexp"
	"strconv"
)

type ContentResult struct {
}

/**
isFilenameValid validates the contents of the filename and returns true if it is ok.
If not, false is returned
*/
func isFilenameValid(fn string) bool {
	matcher := regexp.MustCompile("[;']")
	return !matcher.MatchString(fn)
}

func isOctIdValid(octid string) bool {
	matcher := regexp.MustCompile("\\d+")
	return matcher.MatchString(octid)
}

/**
getIDMapping tries to find an ID mapping record for the given URL, which must contain either a `file` or `octopusid`
parameter
*/
func getIDMapping(ctx context.Context, queryStringParams *map[string]string, config *Config) (*IdMappingRecord, *events.APIGatewayProxyResponse) {
	var idMapping *IdMappingRecord
	var err error

	if fn, haveFn := (*queryStringParams)["file"]; haveFn {
		if isFilenameValid(fn) {
			idMapping, err = IdMappingFromFilebase(ctx, config, fn)
			if err != nil {
				log.Print("ERROR FindContent could not get id mapping: ", err)
				//errorDetail := &ErrorDetail{
				//	ErrorCode:   500,
				//	ErrorString: err.Error(),
				//	FileName:    fn,
				//}
				return nil, MakeResponse(500, GenericErrorBody("Database error"))
			}
		} else {
			//errorDetail := &ErrorDetail{
			//	ErrorCode:   400,
			//	ErrorString: "Invalid filespec",
			//	FileName:    fn,
			//}
			return nil, MakeResponse(400, GenericErrorBody("Invalid filespec"))
		}
	} else if octId, haveOctId := (*queryStringParams)["octopusid"]; haveOctId {
		if isOctIdValid(octId) {
			octIdNum, _ := strconv.ParseInt(octId, 10, 64)
			idMapping, err = IdMappingFromOctid(ctx, config, octIdNum)
		} else {
			return nil, MakeResponse(400, GenericErrorBody("Invalid octid"))
		}
	} else {
		return nil, MakeResponse(400, GenericErrorBody("No search"))
	}

	return idMapping, nil
}

func FindContent(ctx context.Context, queryStringParams *map[string]string, config *Config) (*ContentResult, *events.APIGatewayProxyResponse) {
	//FIXME: no memcache implementation yet, we'll see how necessary it actually is
	idMapping, errResponse := getIDMapping(ctx, queryStringParams, config)
	if errResponse != nil {
		return nil, errResponse
	}

	log.Printf("DEBUGGING got id mapping result %v", idMapping)
	return &ContentResult{}, nil
}
