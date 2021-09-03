package common

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
	"time"
)

func TestNewIdMappingRecord(t *testing.T) {
	sampleData := &map[string]types.AttributeValue{
		"contentId":  &types.AttributeValueMemberS{Value: "somecontent"},
		"filebase":   &types.AttributeValueMemberS{Value: "somefile"},
		"lastupdate": &types.AttributeValueMemberS{Value: "2020-01-02T03:04:05Z"},
		"octopus_id": &types.AttributeValueMemberN{Value: "12345"},
	}

	rec, err := NewIdMappingRecord(sampleData)
	if err != nil {
		t.Errorf("NewIdMappingRecord returned an unexpected error: %s", err)
	} else if rec == nil {
		t.Errorf("NewIdMappingRecord returned a nil value")
	} else {
		if rec.filebase != "somefile" {
			t.Errorf("NewIdMappingRecord returned incorrect %s for filebase", rec.filebase)
		}
		expectedTime := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
		if rec.lastupdate != expectedTime {
			t.Errorf("NewIdMappingRecord returned incorrect time %s, expected %s", rec.lastupdate.String(), expectedTime.String())
		}
		if rec.project != nil {
			t.Error("NewIdMappingRecord returned a value for project when none was set")
		}
		if rec.contentId != "somecontent" {
			t.Errorf("NewIdMappingRecord returned incorrect content id %s", rec.contentId)
		}
		if *rec.octopus_id != 12345 {
			t.Errorf("NewIdMappingRecord returned incorrect octopus id %d", *rec.octopus_id)
		}
	}
}

func TestNewIdMappingRecord_MissingData(t *testing.T) {
	//should error if we are missing content id
	sampleData := &map[string]types.AttributeValue{
		"filebase":   &types.AttributeValueMemberS{Value: "somefile"},
		"lastupdate": &types.AttributeValueMemberS{Value: "2020-01-02T03:04:05Z"},
		"octopus_id": &types.AttributeValueMemberN{Value: "12345"},
	}

	_, err := NewIdMappingRecord(sampleData)
	if err == nil {
		t.Error("NewIdMappingRecord should have errored when there was no content id")
	}
}
