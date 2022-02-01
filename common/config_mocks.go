package common

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ConfigMock struct {
	IdMappingTableVal string
	EncodingsTableVal string
}

func (c *ConfigMock) GetDynamoClient() *dynamodb.Client {
	panic("GetDynamoClient should not be called on the mock")
}

func (c *ConfigMock) IdMappingTable() string {
	return c.IdMappingTableVal
}

func (c *ConfigMock) EncodingsTablePtr() *string {
	copied := c.EncodingsTableVal
	return &copied
}

func (c *ConfigMock) MimeEquivalentsTablePtr() *string {
	return aws.String("mime-equivalents")
}
