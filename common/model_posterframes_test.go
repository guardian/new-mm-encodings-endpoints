package common

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
)

func TestPosterFrameFromDynamo(t *testing.T) {
	rec := &RawDynamoRecord{
		"posterid":   &types.AttributeValueMemberN{Value: "11111"},
		"encodingid": &types.AttributeValueMemberN{Value: "22222"},
		"contentid":  &types.AttributeValueMemberN{Value: "33333"},
		"poster_url": &types.AttributeValueMemberS{Value: "https://some/poster/url"},
		"mime_type":  &types.AttributeValueMemberS{Value: "image/jpeg"},
	}

	result, err := PosterFrameFromDynamo(rec)
	if err != nil {
		t.Errorf("Got unexpected error from PosterFrameFromDynamo: %s", err)
		t.FailNow()
	}

	if result.PosterId != 11111 {
		t.Errorf("Got unexpected result for PosterId: %d", result.PosterId)
	}
	if result.EncodingId != 22222 {
		t.Errorf("Got unexpected result for EncodingId: %d", result.EncodingId)
	}
	if result.ContentId != 33333 {
		t.Errorf("Got unexpected result for ContentId: %d", result.ContentId)
	}
	if result.PosterUrl != "https://some/poster/url" {
		t.Errorf("Got unexpected result for poster url: %s", result.PosterUrl)
	}
	if result.MimeType != "image/jpeg" {
		t.Errorf("Got unexpected result for mime type: %s", result.MimeType)
	}
}
