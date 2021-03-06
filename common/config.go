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

type ConfigImpl struct {
	DyanmoContentTable            string
	idMappingTable                string
	MimeEquivalentsTable          string
	PosterFramesTable             string
	MemcacheHost                  string
	MemcachePort                  int16
	MemcacheExpirySeconds         int16
	MemcacheNotfoundExpirySeconds int16
	awsClientsConfig              aws.Config
	ddbClient                     *dynamodb.Client
}

/*
Config is exposed as an interface to allow for mocking
*/
type Config interface {
	GetDynamoClient() *dynamodb.Client
	IdMappingTable() string
	EncodingsTablePtr() *string
	MimeEquivalentsTablePtr() *string
}

/*
NewConfig initiates a new Config object with values set from default environment variables
*/
func NewConfig() (Config, error) {
	awscfg, awsErr := awsconfig.LoadDefaultConfig(context.Background())
	if awsErr != nil {
		return nil, awsErr
	}

	basicConfig := &ConfigImpl{
		os.Getenv("ENCODINGS_TABLE"),
		os.Getenv("ID_MAPPING_TABLE"),
		os.Getenv("MIME_EQUIVALENTS_TABLE"),
		os.Getenv("POSTER_FRAMES_TABLE"),
		os.Getenv("MEMCACHE_HOST"),
		11211,
		240,
		10,
		awscfg,
		dynamodb.NewFromConfig(awscfg),
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
	if basicConfig.idMappingTable == "" {
		return nil, errors.New("ID_MAPPING_TABLE not set")
	}

	return basicConfig, nil
}

func (c *ConfigImpl) GetDynamoClient() *dynamodb.Client {
	return c.ddbClient
}

func (c *ConfigImpl) IdMappingTable() string {
	return c.idMappingTable
}

func (c *ConfigImpl) EncodingsTablePtr() *string {
	return aws.String(c.DyanmoContentTable)
}

func (c *ConfigImpl) MimeEquivalentsTablePtr() *string {
	return aws.String(c.MimeEquivalentsTable)
}
