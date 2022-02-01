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
	lambda.Start(HandleEvent)
}
