package common

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DyanmoContentTable            string
	IdMappingTable                string
	MemcacheHost                  string
	MemcachePort                  int16
	MemcacheExpirySeconds         int16
	MemcacheNotfoundExpirySeconds int16
	awsClientsConfig              aws.Config
}

/**
NewConfig initiates a new Config object with values set from default environment variables
*/
func NewConfig() (*Config, error) {
	awscfg, awsErr := awsconfig.LoadDefaultConfig(context.Background())
	if awsErr != nil {
		return nil, awsErr
	}

	basicConfig := &Config{
		os.Getenv("CONTENT_TABLE_NAME"),
		os.Getenv("ID_MAPPING_TABLE"),
		os.Getenv("MEMCACHE_HOST"),
		11211,
		240,
		10,
		awscfg,
	}

	if os.Getenv("MEMCACHE_PORT") != "" {
		maybeNewPort, err := strconv.ParseInt(os.Getenv("MEMCACHE_PORT"), 10, 16)
		if err != nil {
			log.Printf("ERROR NewConfig MEMCACHE_PORT is not a valid number: %s", err)
			return nil, errors.New("MEMCACHE_PORT not valid")
		}
		basicConfig.MemcachePort = int16(maybeNewPort)
	}

	if basicConfig.DyanmoContentTable == "" {
		return nil, errors.New("CONTENT_TABLE_NAME is not set")
	}
	if basicConfig.IdMappingTable == "" {
		return nil, errors.New("ID_MAPPING_TABLE not set")
	}
	if basicConfig.MemcacheHost == "" {
		return nil, errors.New("MEMCACHE_HOST is not set")
	}
	return basicConfig, nil
}

func (c *Config) GetDynamoClient() *dynamodb.Client {
	return dynamodb.NewFromConfig(c.awsClientsConfig)
}
