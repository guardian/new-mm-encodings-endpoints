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

var DefaultHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Methods":     "GET, OPTIONS",
	"Access-Control-Allow-Headers":     "*",
	"Access-Control-Allow-Credentials": "false",
	"Access-Control-Max-Age":           "3600",
}

func getDefaultHeaders() map[string]string {
	newHeaders := make(map[string]string, len(DefaultHeaders))
	for k, v := range DefaultHeaders {
		newHeaders[k] = v
	}
	return newHeaders
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
MakeResponseJson takes a JSON-serializable object and returns a pointer to an APIGatewayProxyResponse that can be
passed straight back to API Gateway.

Arguments:
- responseCode - HTTP response code to return to the caller
- contentBody - a JSON-serializable object.  This is serialized to json and the resulting byte stream used as the
response body.  If the serialization fails for any reason then an error is logged and a generic 500 error response is returned.
Returns:
- Pointer to an APIGatewayProxyResponse object
*/
func MakeResponseJson(responseCode int, contentBody interface{}) *events.APIGatewayProxyResponse {
	var jsonBytes []byte
	var stringContent string
	if contentBody != nil {
		var marshalErr error
		jsonBytes, marshalErr = json.Marshal(contentBody)
		if marshalErr != nil {
			log.Printf("ERROR MakeResponseJson could not marshal %v into json: %s", contentBody, marshalErr)
			return MakeResponseJson(500, GenericErrorBody("invalid output content"))
		}
		stringContent = string(jsonBytes)
	} else {
		stringContent = ""
	}

	return MakeResponseRaw(responseCode, &stringContent, "application/json")
}

/*
MakeResponseRaw takes a string and uses it as the content body in a response. It returns a pointer to APIGatewayProxyResponse
that can be passed directly back to the runtime.

Arguments:

- responseCode - HTTP status code

- contentBodyString - the content body to use, or "" for an empty response

- contentType - MIME type to indicate what the content is
*/
func MakeResponseRaw(responseCode int, contentBodyString *string, contentType string) *events.APIGatewayProxyResponse {
	headers := getDefaultHeaders()
	contentLength := len(*contentBodyString)
	if contentLength != 0 {
		headers["Content-Type"] = contentType
		headers["Content-Length"] = fmt.Sprintf("%d", contentLength)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: responseCode,
		Headers:    headers,
		Body:       *contentBodyString,
	}
}

/*
MakeResponseRedirect takes a URL and uses it in the location header in a response. It returns a pointer to APIGatewayProxyResponse
that can be passed directly back to the runtime.

Argument:

- uRL - URL to use in the location header
*/
func MakeResponseRedirect(uRL string) *events.APIGatewayProxyResponse {
	headers := getDefaultHeaders()
	headers["Location"] = fmt.Sprintf("%s", uRL)

	return &events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers:    headers,
	}
}
