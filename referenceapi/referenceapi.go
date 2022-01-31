package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
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

	foundContent, errResponse := common.FindContent(ctx, &event.QueryStringParameters, ops, config)
	if errResponse != nil {
		switch errResponse.StatusCode {
		case 404:
			return common.MakeResponseRaw(404, aws.String("No content found.\n"), "text/plain"), nil
		default:
			return errResponse, nil
		}
	}

	if _, ok := (event.QueryStringParameters)["poster"]; ok {
		if foundContent.PosterURL != "" {
			return common.MakeResponseRaw(200, &foundContent.PosterURL, "text/plain"), nil
		} else {
			return common.MakeResponseRaw(404, aws.String("No poster URL found"), "text/plain"), nil
		}
	}

	return common.MakeResponseRaw(200, &foundContent.Url, "text/plain"), nil
}

func main() {
	lambda.Start(HandleEvent)
}
