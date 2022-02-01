package main

import (
	"context"
	"fmt"
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
This script looks up a video in the interactivepublisher database and returns an HTML video tag with the URL of the video in
*/

func HandleEvent(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	foundContent, errResponse := common.FindContent(ctx, &event.QueryStringParameters, ops, config, mimeEquivelentsCache)
	if errResponse != nil {
		switch errResponse.StatusCode {
		case 404:
			return common.MakeResponseRaw(404, aws.String("No content found.\n"), "text/plain;charset=UTF-8"), nil
		default:
			return errResponse, nil
		}
	}

	extraArguments := ""
	if _, hasNoControls := (event.QueryStringParameters)["nocontrols"]; hasNoControls == false {
		extraArguments = extraArguments + " controls"
	}
	if _, hasAutoPlay := (event.QueryStringParameters)["autoplay"]; hasAutoPlay {
		extraArguments = extraArguments + " autoplay"
	}
	if _, hasLoop := (event.QueryStringParameters)["loop"]; hasLoop {
		extraArguments = extraArguments + " loop"
	}

	hTMLToReturn := "<video preload='auto' id='video_" + fmt.Sprint(foundContent.OctopusId) + "' poster='" + foundContent.PosterURL + "'" + extraArguments + ">\n\t<source src='" + foundContent.Url + "' type='" + foundContent.Format + "'>\n</video>\n"
	return common.MakeResponseRaw(200, &hTMLToReturn, "text/html;charset=UTF-8"), nil
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
