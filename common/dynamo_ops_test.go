package common

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"os"
)

func testAwsConfig() (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:8080",
				SigningRegion: "us-east-1",
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	return config.LoadDefaultConfig(context.Background(), config.WithEndpointResolverWithOptions(customResolver))
}

func setupSomeEncodings(client *dynamodb.Client) {
	createTableRq := &dynamodb.CreateTableInput{
		AttributeDefinitions:   nil,
		KeySchema:              nil,
		TableName:              nil,
		BillingMode:            "",
		GlobalSecondaryIndexes: nil,
		LocalSecondaryIndexes:  nil,
		ProvisionedThroughput:  nil,
		SSESpecification:       nil,
		StreamSpecification:    nil,
		TableClass:             "",
		Tags:                   nil,
	}
	client.CreateTable(context.Background(), createTableRq)
}
func init() {
	if os.Getenv("RUN_DDB_TESTS") != "" {
		cfg, err := testAwsConfig()
		if err != nil {
			log.Fatal("Could not set up AWS config for testing: ", err)
		}

		client := dynamodb.NewFromConfig(cfg)
	} else {
		println("INFO Not performing dynamodb-local test as RUN_DDB_TESTS is not set")
	}
}
