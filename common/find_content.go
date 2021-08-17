package common

import (
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/url"
	"regexp"
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
func getIDMapping(requestUri *url.URL, config *Config) (*IdMappingRecord, *events.APIGatewayProxyResponse) {
	var idMapping *IdMappingRecord
	var err error
	fn := requestUri.Query().Get("file")
	if fn != "" {
		if isFilenameValid(fn) {
			idMapping, err = IdMappingFromFilebase(config, fn)
			if err != nil {
				log.Print("ERROR FindContent could not get id mapping: ", err)
				errorDetail := &ErrorDetail{
					ErrorCode:   500,
					ErrorString: err.Error(),
					FileName:    fn,
					QueryUrl:    requestUri.String(),
				}
				return nil, MakeResponse(500, GenericErrorBody("Database error"))
			}
		} else {
			errorDetail := &ErrorDetail{
				ErrorCode:   400,
				ErrorString: "Invalid filespec",
				FileName:    fn,
				QueryUrl:    requestUri.String(),
			}
			return nil, MakeResponse(400, GenericErrorBody("Invalid filespec"))
		}
	}

	octId := requestUri.Query().Get("octopusid")
	if fn == "" && octId != "" {
		if isOctIdValid(octId) {

		} else {
			return nil, MakeResponse(400, GenericErrorBody("Invalid octid"))
		}
	}

	if fn == "" && octId == "" {

	}

	return idMapping, nil
}
func FindContent(requestUri *url.URL, config *Config) (*ContentResult, *events.APIGatewayProxyResponse) {
	//FIXME: no memcache implementation yet, we'll see how necessary it actually is

}
