package common

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
	"time"
)

/*
Tests that a valid record can be decoded
*/
func TestEncodingFromDynamo(t *testing.T) {
	testrecord := &RawDynamoRecord{
		"encodingid":   &types.AttributeValueMemberN{Value: "1234"},
		"contentid":    &types.AttributeValueMemberN{Value: "2345"},
		"url":          &types.AttributeValueMemberS{Value: "http://some/encoding/urk"},
		"format":       &types.AttributeValueMemberS{Value: "video/mp4"},
		"mobile":       &types.AttributeValueMemberBOOL{Value: false},
		"multirate":    &types.AttributeValueMemberBOOL{Value: false},
		"vcodec":       &types.AttributeValueMemberS{Value: "h264"},
		"acodec":       &types.AttributeValueMemberS{Value: "aac"},
		"lastupdate":   &types.AttributeValueMemberS{Value: "2016-02-03T04:05:06Z"},
		"frame_width":  &types.AttributeValueMemberN{Value: "1280"},
		"frame_height": &types.AttributeValueMemberN{Value: "720"},
		"duration":     &types.AttributeValueMemberN{Value: "123.456"},
		"file_size":    &types.AttributeValueMemberN{Value: "567890123"},
		"fcs_id":       &types.AttributeValueMemberS{Value: "9999998"},
		"octopus_id":   &types.AttributeValueMemberN{Value: "55543"},
		"aspect":       &types.AttributeValueMemberS{Value: "16:9"},
	}

	result, err := EncodingFromDynamo(testrecord)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expectedUpdateTime, _ := time.Parse(time.RFC3339, "2016-02-03T04:05:06Z")

	if result.EncodingId != 1234 {
		t.Errorf("Unexpected encoding id %d", result.EncodingId)
	}
	if result.ContentId != 2345 {
		t.Errorf("Unexpected content id %d", result.ContentId)
	}
	if result.Url != "http://some/encoding/urk" {
		t.Errorf("Unexpected url '%s'", result.Url)
	}
	if result.Format != "video/mp4" {
		t.Errorf("Unexpected format '%s'", result.Format)
	}
	if result.Mobile != false {
		t.Errorf("Unexpected mobile result: %v", result.Mobile)
	}
	if result.Multirate != false {
		t.Errorf("Unexpected multirate result %v", result.Multirate)
	}
	if result.VCodec != "h264" {
		t.Errorf("Unexpected vcodec %s", result.VCodec)
	}
	if result.ACodec != "aac" {
		t.Errorf("Unexpected acodec %s", result.ACodec)
	}
	if result.LastUpdate != expectedUpdateTime {
		t.Errorf("Unexpected update time %s", result.LastUpdate.Format(time.RFC3339))
	}
	if result.FrameWidth != 1280 {
		t.Errorf("Unexpected frame width %d", result.FrameWidth)
	}
	if result.FrameHeight != 720 {
		t.Errorf("Unexpected frame height %d", result.FrameHeight)
	}
	if result.Duration != 123.456 {
		t.Errorf("Unexpected duration %f", result.Duration)
	}
	if result.FileSize != 567890123 {
		t.Errorf("Unexpected file size %d", result.FileSize)
	}
	if result.FCSID != "9999998" {
		t.Errorf("Unexpected FCS ID %s", result.FCSID)
	}
	if result.OctopusId != 55543 {
		t.Errorf("Unexpected Octopus ID")
	}
	if result.Aspect != "16:9" {
		t.Errorf("Unexpected aspect ratio %s", result.Aspect)
	}
}

/*
test that a record missing optional fields can still be decoded
*/
func TestEncodingFromDynamoMissingOptional(t *testing.T) {
	testrecord := &RawDynamoRecord{
		"encodingid":   &types.AttributeValueMemberN{Value: "1234"},
		"contentid":    &types.AttributeValueMemberN{Value: "2345"},
		"url":          &types.AttributeValueMemberS{Value: "http://some/encoding/urk"},
		"format":       &types.AttributeValueMemberS{Value: "video/mp4"},
		"mobile":       &types.AttributeValueMemberBOOL{Value: false},
		"multirate":    &types.AttributeValueMemberBOOL{Value: false},
		"lastupdate":   &types.AttributeValueMemberS{Value: "2016-02-03T04:05:06Z"},
		"frame_width":  &types.AttributeValueMemberN{Value: "1280"},
		"frame_height": &types.AttributeValueMemberN{Value: "720"},
		"duration":     &types.AttributeValueMemberN{Value: "123.456"},
		"file_size":    &types.AttributeValueMemberN{Value: "567890123"},
		"fcs_id":       &types.AttributeValueMemberS{Value: "9999998"},
		"octopus_id":   &types.AttributeValueMemberN{Value: "55543"},
		"aspect":       &types.AttributeValueMemberS{Value: "16:9"},
	}

	result, err := EncodingFromDynamo(testrecord)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expectedUpdateTime, _ := time.Parse(time.RFC3339, "2016-02-03T04:05:06Z")

	if result.EncodingId != 1234 {
		t.Errorf("Unexpected encoding id %d", result.EncodingId)
	}
	if result.ContentId != 2345 {
		t.Errorf("Unexpected content id %d", result.ContentId)
	}
	if result.Url != "http://some/encoding/urk" {
		t.Errorf("Unexpected url '%s'", result.Url)
	}
	if result.Format != "video/mp4" {
		t.Errorf("Unexpected format '%s'", result.Format)
	}
	if result.Mobile != false {
		t.Errorf("Unexpected mobile result: %v", result.Mobile)
	}
	if result.Multirate != false {
		t.Errorf("Unexpected multirate result %v", result.Multirate)
	}
	if result.LastUpdate != expectedUpdateTime {
		t.Errorf("Unexpected update time %s", result.LastUpdate.Format(time.RFC3339))
	}
	if result.FrameWidth != 1280 {
		t.Errorf("Unexpected frame width %d", result.FrameWidth)
	}
	if result.FrameHeight != 720 {
		t.Errorf("Unexpected frame height %d", result.FrameHeight)
	}
	if result.Duration != 123.456 {
		t.Errorf("Unexpected duration %f", result.Duration)
	}
	if result.FileSize != 567890123 {
		t.Errorf("Unexpected file size %d", result.FileSize)
	}
	if result.FCSID != "9999998" {
		t.Errorf("Unexpected FCS ID %s", result.FCSID)
	}
	if result.OctopusId != 55543 {
		t.Errorf("Unexpected Octopus ID")
	}
	if result.Aspect != "16:9" {
		t.Errorf("Unexpected aspect ratio %s", result.Aspect)
	}
}

/*
Tests that we fail if mandatory fields are missing
*/
func TestEncodingFromDynamoMissingMandatory(t *testing.T) {
	testrecord := &RawDynamoRecord{
		"encodingid":   &types.AttributeValueMemberN{Value: "1234"},
		"contentid":    &types.AttributeValueMemberN{Value: "2345"},
		"url":          &types.AttributeValueMemberS{Value: "http://some/encoding/urk"},
		"mobile":       &types.AttributeValueMemberBOOL{Value: false},
		"multirate":    &types.AttributeValueMemberBOOL{Value: false},
		"vcodec":       &types.AttributeValueMemberS{Value: "h264"},
		"acodec":       &types.AttributeValueMemberS{Value: "aac"},
		"lastupdate":   &types.AttributeValueMemberS{Value: "2016-02-03T04:05:06Z"},
		"frame_width":  &types.AttributeValueMemberN{Value: "1280"},
		"frame_height": &types.AttributeValueMemberN{Value: "720"},
		"duration":     &types.AttributeValueMemberN{Value: "123.456"},
		"file_size":    &types.AttributeValueMemberN{Value: "567890123"},
		"fcs_id":       &types.AttributeValueMemberS{Value: "9999998"},
		"octopus_id":   &types.AttributeValueMemberN{Value: "55543"},
		"aspect":       &types.AttributeValueMemberS{Value: "16:9"},
	}

	_, err := EncodingFromDynamo(testrecord)
	if err == nil {
		t.Error("Expected an error on an invalid record but got nothing")
		t.FailNow()
	}

	if err.Error() != "the given record is not a ContentResult" {
		t.Errorf("Unexpected error string on invalid record: %s", err.Error())
	}
}
