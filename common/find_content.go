package common

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"regexp"
	"strconv"
)

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

/*
getFCSId returns the FCS ID for a given contentId.

The FCS id uniquely identifies the version (as opposed to octopus_id uniquely identifies the title which can have
multiple versions.
Versions can have subtly different bitrates AND arrive at different times, so just searching versions with a sort order
can return old results no matter what.
So, the first step is to find the most recent FCS ID and then search with that
Some entries may not have FCS IDs, and if uncaught this leads to all such entries being treated as the same title.
So, we iterate across them all and get the first non-empty one. If no ids are found then we must fall back to the
old behaviour (step 3)
*/
func getFCSId(ctx context.Context, ops DynamoDbOps, contentId int32) (*string, error) {
	results, err := ops.QueryFCSIdForContentId(ctx, contentId)

	if err != nil {
		return nil, err
	}
	for _, r := range *results {
		if r != "" {
			finalResult := r
			return &finalResult, nil
		}
	}
	return nil, nil
}

/**
getIDMapping tries to find an ID mapping record for the given URL, which must contain either a `file` or `octopusid`
parameter
*/
func getIDMapping(ctx context.Context, queryStringParams *map[string]string, config Config) (*IdMappingRecord, *events.APIGatewayProxyResponse) {
	var idMapping *IdMappingRecord
	var err error

	if fn, haveFn := (*queryStringParams)["file"]; haveFn {
		if isFilenameValid(fn) {
			idMapping, err = IdMappingFromFilebase(ctx, config, fn)
			if err != nil {
				log.Print("ERROR FindContent could not get id mapping: ", err)
				return nil, MakeResponse(500, GenericErrorBody("Database error"))
			}
		} else {
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

func FindContent(ctx context.Context, queryStringParams *map[string]string, config Config) (*ContentResult, *events.APIGatewayProxyResponse) {
	ops := NewDynamoDbOps(config)
	//FIXME: no memcache implementation yet, we'll see how necessary it actually is
	idMapping, errResponse := getIDMapping(ctx, queryStringParams, config)
	if errResponse != nil {
		return nil, errResponse
	}

	log.Printf("DEBUGGING got id mapping result %v", idMapping)
	if idMapping != nil { //we got a result from idmapping
		fcsId, err := getFCSId(ctx, ops, idMapping.contentId)
		if err != nil {
			return nil, MakeResponse(500, GenericErrorBody("Database error"))
		}
		if fcsId != nil {
			log.Printf("DEBUGGING got FCS ID %s", *fcsId)
		} else {
			log.Printf("DEBUGGING did not find an FCS ID")
		}
		return &ContentResult{}, nil
	} else { //fall back to direct query
		return nil, MakeResponse(500, GenericErrorBody("Not implemented yet"))
	}
}
