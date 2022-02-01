package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/guardian/new-encodings-endpoints/common"
	"log"
)

var ops common.DynamoDbOps
var config common.Config
var mimeEquivelentsCache common.MimeEquivalentsCache

/*
This script looks up a video in the interactivepublisher database and returns a URL, if it can be found, in a location header
*/

func HandleEvent(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	foundContent, errResponse := common.FindContent(ctx, &event.QueryStringParameters, ops, config, mimeEquivelentsCache)
	if errResponse != nil {
		switch errResponse.StatusCode {
		case 404:
			return common.MakeResponseRaw(404, aws.String("No content found.\n"), "text/plain"), nil
		default:
			return errResponse, nil
		}
	}

	if _, havePoster := (event.QueryStringParameters)["poster"]; havePoster {
		if foundContent.PosterURL != "" {
			return common.MakeResponseRedirect(foundContent.PosterURL), nil
		} else {
			return common.MakeResponseRaw(404, aws.String("No poster URL found"), "text/plain"), nil
		}
	}

	return common.MakeResponseRedirect(foundContent.Url), nil
}

func main() {
	var err error
	config, err = common.NewConfig()
	if err != nil {
		log.Printf("ERROR Could not initialise config: %s", err)
		panic("could not initialise config")
	}

	ops = common.NewDynamoDbOps(config)
	mimeEquivelentsCache, err = common.NewMimeEquivalentsCache(context.Background(), ops)
	if err != nil {
		log.Printf("ERROR Could not initialise mime equivalents: %s", err)
		panic("could not initialise MIME equivalents")
	}

	lambda.Start(HandleEvent)
}
