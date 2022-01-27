package common

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"regexp"
	"strconv"
	"time"
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
	matcher := regexp.MustCompile("^\\d+$")
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
	if results == nil {
		return nil, nil
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
func getIDMapping(ctx context.Context, queryStringParams *map[string]string, ops DynamoDbOps, config Config) (*IdMappingRecord, *events.APIGatewayProxyResponse) {
	var idMapping *IdMappingRecord
	var err error

	if fn, haveFn := (*queryStringParams)["file"]; haveFn {
		if isFilenameValid(fn) {
			idMapping, err = ops.QueryIdMappings(ctx, IdMappingIndexFilebase, IdMappingKeyfieldFilebase, fn)
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
			idMapping, err = ops.QueryIdMappings(ctx, IdMappingIndexOctid, IdMappingKeyfieldOctid, octIdNum)
		} else {
			return nil, MakeResponse(400, GenericErrorBody("Invalid octid"))
		}
	} else {
		return nil, MakeResponse(400, GenericErrorBody("No search"))
	}

	return idMapping, nil
}

/*
FindContent is the main entry point to the common logic for all the endpoints. It takes in the query parameters and tries to
find the best match for them, returning this as a pointer to ContentResult.

Arguments:
- ctx - context that can be used to cancel the operation, passed in from lambda functions
- queryStringParams - pointer to a string-string map representing the query parameters from the URL
- ops - a DynamoDbOps object that abstracts the actual Dynamo operations for mocking
- config - a Config object that encapsulates the runtime configuration
Returns:
- a pointer to ContentResult on success
- a pointer to APIGatewayProxyResponse on error. This can be passed back directly to the runtime.
*/
func FindContent(ctx context.Context, queryStringParams *map[string]string, ops DynamoDbOps, config Config) (*ContentResult, *events.APIGatewayProxyResponse) {
	//FIXME: no memcache implementation yet, we'll see how necessary it actually is
	idMapping, errResponse := getIDMapping(ctx, queryStringParams, ops, config)
	if errResponse != nil {
		return nil, errResponse
	}

	var contentToFilter []*Encoding
	log.Printf("DEBUGGING got id mapping result %v", idMapping)
	if idMapping == nil { //nothing in idmapping => does not exist
		return nil, MakeResponse(404, GenericErrorBody("Content not found"))
	}
	var err error

	fcsId, err := getFCSId(ctx, ops, idMapping.contentId)
	if err != nil {
		return nil, MakeResponse(500, GenericErrorBody("Database error"))
	}

	if fcsId != nil {
		log.Printf("DEBUGGING got FCS ID %s", *fcsId)
		contentToFilter, err = ops.QueryEncodingsForFCSId(ctx, *fcsId)
		if err != nil {
			log.Printf("ERROR Could not query encodings: %s", err)
			return nil, MakeResponse(500, GenericErrorBody("Database error"))
		}
	}

	if contentToFilter == nil { //we didn't get any results yet
		log.Print("INFO No content from primary search, falling back to secondary")
		_, haveAllowOld := (*queryStringParams)["allow_old"]
		var maybeSince *time.Time
		if !haveAllowOld {
			maybeSince = &idMapping.lastupdate
			log.Printf("INFO allow_old not set, only looking for results since %s", maybeSince.Format(time.RFC3339))
		}
		contentToFilter, err = ops.QueryEncodingsForContentId(ctx, idMapping.contentId, maybeSince)
		if err != nil {
			log.Printf("ERROR Could not query encodings: %s", err)
			return nil, MakeResponse(500, GenericErrorBody("Database error"))
		}
	}

	for _, c := range contentToFilter {
		log.Printf("INFO Got record %v", *c)
	}

	var format = ""
	if val, ok := (*queryStringParams)["format"]; ok {
		format = val
	}

	var need_mobile = false
	if val, ok := (*queryStringParams)["need_mobile"]; ok {
		if val == "true" {
			need_mobile = true
		}
	}

	var minbitrate int32 = 0
	if val, ok := (*queryStringParams)["minbitrate"]; ok {
		var parseIntOutput, _ = strconv.ParseInt(val, 10, 32)
		minbitrate = int32(parseIntOutput)
	}

	var maxbitrate int32 = 0
	if val, ok := (*queryStringParams)["maxbitrate"]; ok {
		var parseIntOutput, _ = strconv.ParseInt(val, 10, 32)
		maxbitrate = int32(parseIntOutput)
	}

	var minheight int32 = 0
	if val, ok := (*queryStringParams)["minheight"]; ok {
		var parseIntOutput, _ = strconv.ParseInt(val, 10, 32)
		minheight = int32(parseIntOutput)
	}

	var maxheight int32 = 0
	if val, ok := (*queryStringParams)["maxheight"]; ok {
		var parseIntOutput, _ = strconv.ParseInt(val, 10, 32)
		maxheight = int32(parseIntOutput)
	}

	var minwidth int32 = 0
	if val, ok := (*queryStringParams)["minwidth"]; ok {
		var parseIntOutput, _ = strconv.ParseInt(val, 10, 32)
		minwidth = int32(parseIntOutput)
	}

	var maxwidth int32 = 0
	if val, ok := (*queryStringParams)["maxwidth"]; ok {
		var parseIntOutput, _ = strconv.ParseInt(val, 10, 32)
		maxwidth = int32(parseIntOutput)
	}

	if len(contentToFilter) > 0 {
		filteredContent := ContentFilter(contentToFilter, format, need_mobile, minbitrate, maxbitrate, minheight, maxheight, minwidth, maxwidth)
		return filteredContent, nil
	} else {
		return nil, MakeResponse(404, GenericErrorBody("No encodings matching your request"))
	}
}
