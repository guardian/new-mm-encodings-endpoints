package common

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
)

func TestMimeEquivalentFromDynamo(t *testing.T) {
	rec := &RawDynamoRecord{
		"real_name":       &types.AttributeValueMemberS{Value: "video/m3u8"},
		"mime_equivalent": &types.AttributeValueMemberS{Value: "application/x-m3u8"},
	}

	result, err := MimeEquivalentFromDynamo(rec)
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
		t.FailNow()
	}

	if result.MimeEquivalent != "application/x-m3u8" {
		t.Errorf("Got incorrect MimeEquivalent '%s'", result.MimeEquivalent)
	}

	if result.RealName != "video/m3u8" {
		t.Errorf("Got incorrect RealName '%s'", result.RealName)
	}
}
