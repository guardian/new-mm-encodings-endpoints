package main

//This function returns a permissive CORS header in response to an OPTIONS preflight request
import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/guardian/new-encodings-endpoints/common"
)

func HandleEvent(ctx context.Context, event *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	//the standard CORS headers are all included in the `MakeResponse` output
	return common.MakeResponseRaw(200, aws.String(""), ""), nil
}

func main() {
	lambda.Start(HandleEvent)
}
