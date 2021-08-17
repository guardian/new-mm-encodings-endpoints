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

func GenericErrorBody(msg string) map[string]string {
	return map[string]string{
		"status": "error",
		"detail": msg,
	}
}

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
