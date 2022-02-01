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
This script looks up a video in the interactivepublisher database and returns a URL, if it can be found, in a location header
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
			return common.MakeResponseRedirect(foundContent.PosterURL), nil
		} else {
			return common.MakeResponseRaw(404, aws.String("No poster URL found"), "text/plain"), nil
		}
	}

	return common.MakeResponseRedirect(foundContent.Url), nil
}

func main() {
	lambda.Start(HandleEvent)
}
