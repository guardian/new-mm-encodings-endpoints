package common

import (
	"context"
	"strings"
	"testing"
	"time"
)

/*
FindContent should return a 400 error if we don't have either recognised search parameter
*/
func TestFindContentInvalidParams(t *testing.T) {
	fakeParams := map[string]string{"foo": "bar"}
	tim, _ := time.Parse(time.RFC3339, time.RFC3339)
	ops := &DynamoOpsMock{
		IdMappingResult: IdMappingRecord{
			contentId:  2222,
			filebase:   "something",
			project:    nil,
			lastupdate: tim,
			octopus_id: nil,
		},
	}
	config := &ConfigMock{
		IdMappingTableVal: "id-mapping-table",
		EncodingsTableVal: "encodings-table",
	}

	content, errResponse := FindContent(context.Background(), &fakeParams, ops, config)
	if content != nil {
		t.Error("FindContent returned content result for an invalid query")
	}
	if errResponse == nil {
		t.Error("FindContent did not return an error response for an invalid query")
	} else {
		if errResponse.StatusCode != 400 {
			t.Errorf("FindContent returned wrong status code for invalid query, got %d wanted 400", errResponse.StatusCode)
		}
		if !strings.Contains(errResponse.Body, "No search") {
			t.Errorf("FindContent returned content body '%s' did not include the expected error string", errResponse.Body)
		}
	}
}

/*
FindContent should return a 400 error if we have a suspicious filename
*/
func TestFindContentInvalidFilename(t *testing.T) {
	fakeParams := map[string]string{"file": "somethingsomething; drop table everything; haha"}
	tim, _ := time.Parse(time.RFC3339, time.RFC3339)
	ops := &DynamoOpsMock{
		IdMappingResult: IdMappingRecord{
			contentId:  2222,
			filebase:   "something",
			project:    nil,
			lastupdate: tim,
			octopus_id: nil,
		},
	}
	config := &ConfigMock{
		IdMappingTableVal: "id-mapping-table",
		EncodingsTableVal: "encodings-table",
	}

	content, errResponse := FindContent(context.Background(), &fakeParams, ops, config)
	if content != nil {
		t.Error("FindContent returned content result for an invalid filebase")
	}
	if errResponse == nil {
		t.Error("FindContent did not return an error response for an invalid filebase")
	} else {
		if errResponse.StatusCode != 400 {
			t.Errorf("FindContent returned wrong status code for invalid filebase, got %d wanted 400", errResponse.StatusCode)
		}
		if !strings.Contains(errResponse.Body, "Invalid filespec") {
			t.Errorf("FindContent returned content body '%s' did not include the expected error string", errResponse.Body)
		}
	}
}

/*
FindContent should return a 400 error if we have a suspicious filename
*/
func TestFindContentValidFilename(t *testing.T) {
	fakeParams := map[string]string{"file": "mygreatvideo"}
	tim, _ := time.Parse(time.RFC3339, time.RFC3339)
	ops := &DynamoOpsMock{
		IdMappingResult: IdMappingRecord{
			contentId:  2222,
			filebase:   "something",
			project:    nil,
			lastupdate: tim,
			octopus_id: nil,
		},
		FCSIdForContentIdResults: &[]string{"KP-12345", "KP-12345"},
	}

	config := &ConfigMock{
		IdMappingTableVal: "id-mapping-table",
		EncodingsTableVal: "encodings-table",
	}

	content, errResponse := FindContent(context.Background(), &fakeParams, ops, config)
	if errResponse != nil {
		t.Errorf("FindContent returned an error '%v' for a valid filebase", errResponse)
		t.FailNow()
	}

	if content == nil {
		t.Error("FindContent returned nil result for a valid filebase")
	}
}
