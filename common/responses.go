package common

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

type ErrorDetail struct {
	ErrorCode   int    `json:"error_code"`
	ErrorString string `json:"error_string"`
	FileName    string `json:"file_name"`
	QueryUrl    string `json:"query_url"`
}

/*
GenericErrorBody generates a standard error response suitable for serialization into json
Arguments: `msg` - string to include as the "detail" field in the returned json object
*/
func GenericErrorBody(msg string) map[string]string {
	return map[string]string{
		"status": "error",
		"detail": msg,
	}
}

/*
MakeResponse takes a JSON-serializable object and returns a pointer to an APIGatewayProxyResponse that can be
passed straight back to API Gateway.

Arguments:
- responseCode - HTTP response code to return to the caller
- contentBody - a JSON-serializable object.  This is serialized to json and the resulting byte stream used as the
response body.  If the serialization fails for any reason then an error is logged and a generic 500 error response is returned.
Returns:
- Pointer to an APIGatewayProxyResponse object
*/
func MakeResponse(responseCode int, contentBody interface{}) *events.APIGatewayProxyResponse {
	jsonBytes, marshalErr := json.Marshal(contentBody)
	if marshalErr != nil {
		log.Printf("ERROR MakeResponse could not marshal %v into json: %s", contentBody, marshalErr)
		return MakeResponse(500, GenericErrorBody("invalid output content"))
	} else {
		contentLength := len(jsonBytes)
		return &events.APIGatewayProxyResponse{
			StatusCode: responseCode,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Content-Length":              fmt.Sprintf("%d", contentLength),
				"Access-Control-Allow-Origin": "*",
			},
			Body: string(jsonBytes),
		}
	}
}
