package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/guardian/new-encodings-endpoints/common"
	"log"
)

/*
This script looks up a video in the interactivepublisher database and returns a plaintext url if it can be found
*/

func HandleEvent(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	config, configErr := common.NewConfig()
	if configErr != nil {
		log.Printf("ERROR Could not initialise config: %s", configErr)
		return nil, configErr
	}

	ops := common.NewDynamoDbOps(config)

	_, errResponse := common.FindContent(ctx, &event.QueryStringParameters, ops, config)
	if errResponse != nil {
		return errResponse, nil
	}

	return common.MakeResponse(200, map[string]string{"status": "ok", "detail": "testing"}), nil
}

func main() {
	lambda.Start(HandleEvent)
}
