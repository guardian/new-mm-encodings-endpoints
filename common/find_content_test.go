package common

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
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
FindContent should return a 200 with content if we found something
FIXME: No data is returned yet so we need to update the test to check the return value when it does
*/
func TestFindContentValidFilename(t *testing.T) {
	fakeParams := map[string]string{"file": "mygreatvideo"}
	tim, _ := time.Parse(time.RFC3339, time.RFC3339)
	ops := &DynamoOpsMock{
		IdMappingResult: IdMappingRecord{
			contentId:  2222,
			filebase:   "mygreatvideo",
			project:    nil,
			lastupdate: tim,
			octopus_id: nil,
		},
		FCSIdForContentIdResults: &[]string{"KP-12345", "KP-12346", "KP-12347"},
		EncodingsForFCSIdResults: []*Encoding{
			&Encoding{
				EncodingId:  123,
				ContentId:   111,
				Url:         "https://url/to/content",
				Format:      "mp4",
				Mobile:      false,
				Multirate:   false,
				VCodec:      "h264",
				ACodec:      "aac",
				VBitrate:    12345,
				ABitrate:    128,
				LastUpdate:  tim,
				FrameWidth:  1280,
				FrameHeight: 720,
				Duration:    123.456,
				FileSize:    98765432,
				FCSID:       "KP-12345",
				OctopusId:   34567,
				Aspect:      "16:9",
			},
		},
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

	if ops.FCSIdQueried != "KP-12345" {
		t.Errorf("FindContent queried the wrong FCS ID, expected KP-12345 got '%s'", ops.FCSIdQueried)
	}
	if ops.ContentIdQueried != 0 {
		t.Errorf("FindContent should not have queried by content ID but we got %d", ops.ContentIdQueried)
	}
	if ops.LastContentId != 2222 {
		t.Errorf("FindContent should have queried FCS ID for content ID 2222 but we got %d", ops.LastContentId)
	}
	if ops.IdMappingIndexQueried != IdMappingIndexFilebase {
		t.Errorf("FindContent should have queried id mappings by filebase but we got %s", ops.IdMappingIndexQueried)
	}
	if ops.IdMappingKeyFieldQueried != IdMappingKeyfieldFilebase {
		t.Errorf("FindContent should have queried id mappings on field filebase but we got %s", ops.IdMappingKeyFieldQueried)
	}
	if ops.IdMappingSearchTermQueried != "mygreatvideo" {
		t.Errorf("FindContent queried ID mappings on the wrong term, got %s", ops.IdMappingSearchTermQueried)
	}
}

/*
FindContent should return a 200 with content if we found something
FIXME: No data is returned yet so we need to update the test to check the return value when it does
*/
func TestFindContentValidFilenameNotFound(t *testing.T) {
	fakeParams := map[string]string{"file": "mygreatvideo"}
	tim, _ := time.Parse(time.RFC3339, time.RFC3339)
	ops := &DynamoOpsMock{
		FCSIdForContentIdResults: &[]string{"KP-12345", "KP-12346", "KP-12347"},
		EncodingsForFCSIdResults: []*Encoding{
			&Encoding{
				EncodingId:  123,
				ContentId:   111,
				Url:         "https://url/to/content",
				Format:      "mp4",
				Mobile:      false,
				Multirate:   false,
				VCodec:      "h264",
				ACodec:      "aac",
				VBitrate:    12345,
				ABitrate:    128,
				LastUpdate:  tim,
				FrameWidth:  1280,
				FrameHeight: 720,
				Duration:    123.456,
				FileSize:    98765432,
				FCSID:       "KP-12345",
				OctopusId:   34567,
				Aspect:      "16:9",
			},
		},
	}

	config := &ConfigMock{
		IdMappingTableVal: "id-mapping-table",
		EncodingsTableVal: "encodings-table",
	}

	content, errResponse := FindContent(context.Background(), &fakeParams, ops, config)
	if errResponse == nil {
		t.Errorf("FindContent returned no error for a file not found")
		t.FailNow()
	}

	if errResponse.StatusCode != 404 {
		t.Errorf("FindContent returned error %d for file not found, expected 404", errResponse.StatusCode)
	}

	if !strings.Contains(errResponse.Body, "Content not found") {
		t.Errorf("Did not found the 'Not Found' string in the error body, got %s", errResponse.Body)
	}
	if content != nil {
		t.Error("FindContent returned content for an id not found")
	}

}

/*
FindContent should return a 400 error if we have a suspicious octid
*/
func TestFindContentInvalidOctopusId(t *testing.T) {
	fakeParams := map[string]string{"octopusid": "12345; drop table everything; haha"}
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
		t.Error("FindContent returned content result for an invalid octid")
	}
	if errResponse == nil {
		t.Error("FindContent did not return an error response for an invalid octid")
	} else {
		if errResponse.StatusCode != 400 {
			t.Errorf("FindContent returned wrong status code for invalid octid, got %d wanted 400", errResponse.StatusCode)
		}
		if !strings.Contains(errResponse.Body, "Invalid octid") {
			t.Errorf("FindContent returned content body '%s' did not include the expected error string", errResponse.Body)
		}
	}
}

/*
FindContent should return a 400 error if we have a suspicious octid
*/
func TestFindContentAnotherInvalidOctopusId(t *testing.T) {
	fakeParams := map[string]string{"octopusid": "1234a"}
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
		t.Error("FindContent returned content result for an invalid octid")
	}
	if errResponse == nil {
		t.Error("FindContent did not return an error response for an invalid octid")
	} else {
		if errResponse.StatusCode != 400 {
			t.Errorf("FindContent returned wrong status code for invalid octid, got %d wanted 400", errResponse.StatusCode)
		}
		if !strings.Contains(errResponse.Body, "Invalid octid") {
			t.Errorf("FindContent returned content body '%s' did not include the expected error string", errResponse.Body)
		}
	}
}

/*
FindContent should return a 200 with content if we found something
FIXME: No data is returned yet so we need to update the test to check the return value when it does
*/
func TestFindContentValidOctid(t *testing.T) {
	fakeParams := map[string]string{"octopusid": "123456"}
	tim, _ := time.Parse(time.RFC3339, time.RFC3339)
	ops := &DynamoOpsMock{
		IdMappingResult: IdMappingRecord{
			contentId:  2222,
			filebase:   "mygreatvideo",
			project:    nil,
			lastupdate: tim,
			octopus_id: aws.Int64(123456),
		},
		FCSIdForContentIdResults: &[]string{"KP-12345", "KP-12346", "KP-12347"},
		EncodingsForFCSIdResults: []*Encoding{
			&Encoding{
				EncodingId:  123,
				ContentId:   111,
				Url:         "https://url/to/content",
				Format:      "mp4",
				Mobile:      false,
				Multirate:   false,
				VCodec:      "h264",
				ACodec:      "aac",
				VBitrate:    12345,
				ABitrate:    128,
				LastUpdate:  tim,
				FrameWidth:  1280,
				FrameHeight: 720,
				Duration:    123.456,
				FileSize:    98765432,
				FCSID:       "KP-12345",
				OctopusId:   123456,
				Aspect:      "16:9",
			},
		},
	}

	config := &ConfigMock{
		IdMappingTableVal: "id-mapping-table",
		EncodingsTableVal: "encodings-table",
	}

	content, errResponse := FindContent(context.Background(), &fakeParams, ops, config)
	if errResponse != nil {
		t.Errorf("FindContent returned an error '%v' for a valid octopusid", errResponse)
		t.FailNow()
	}

	if content == nil {
		t.Error("FindContent returned nil result for a valid octopusid")
	}

	if ops.FCSIdQueried != "KP-12345" {
		t.Errorf("FindContent queried the wrong FCS ID, expected KP-12345 got '%s'", ops.FCSIdQueried)
	}
	if ops.ContentIdQueried != 0 {
		t.Errorf("FindContent should not have queried by content ID but we got %d", ops.ContentIdQueried)
	}
	if ops.LastContentId != 2222 {
		t.Errorf("FindContent should have queried FCS ID for content ID 2222 but we got %d", ops.LastContentId)
	}
	if ops.IdMappingIndexQueried != IdMappingIndexOctid {
		t.Errorf("FindContent should have queried id mappings by filebase but we got %s", ops.IdMappingIndexQueried)
	}
	if ops.IdMappingKeyFieldQueried != IdMappingKeyfieldOctid {
		t.Errorf("FindContent should have queried id mappings on field filebase but we got %s", ops.IdMappingKeyFieldQueried)
	}
	if ops.IdMappingSearchTermQueried != int64(123456) {
		t.Errorf("FindContent queried ID mappings on the wrong term, got %d", ops.IdMappingSearchTermQueried)
	}
}

/*
FindContent should fall back to using a more basic type of search if no FCS ID was found
FIXME: No data is returned yet so we need to update the test to check the return value when it does
*/
func TestFindContentFallback(t *testing.T) {
	fakeParams := map[string]string{"octopusid": "123456"}
	tim, _ := time.Parse(time.RFC3339, "2015-01-02T03:04:05Z")
	ops := &DynamoOpsMock{
		IdMappingResult: IdMappingRecord{
			contentId:  2222,
			filebase:   "mygreatvideo",
			project:    nil,
			lastupdate: tim,
			octopus_id: aws.Int64(123456),
		},
		FCSIdForContentIdResults: nil,
		EncodingsForContentIdResults: []*Encoding{
			&Encoding{
				EncodingId:  123,
				ContentId:   111,
				Url:         "https://url/to/content",
				Format:      "mp4",
				Mobile:      false,
				Multirate:   false,
				VCodec:      "h264",
				ACodec:      "aac",
				VBitrate:    12345,
				ABitrate:    128,
				LastUpdate:  tim,
				FrameWidth:  1280,
				FrameHeight: 720,
				Duration:    123.456,
				FileSize:    98765432,
				FCSID:       "KP-12345",
				OctopusId:   123456,
				Aspect:      "16:9",
			},
		},
	}

	config := &ConfigMock{
		IdMappingTableVal: "id-mapping-table",
		EncodingsTableVal: "encodings-table",
	}

	content, errResponse := FindContent(context.Background(), &fakeParams, ops, config)
	if errResponse != nil {
		t.Errorf("FindContent returned an error '%v' for a valid octopusid", errResponse)
		t.FailNow()
	}

	if content == nil {
		t.Error("FindContent returned nil result for a valid octopusid")
	}

	if ops.FCSIdQueried != "" {
		t.Errorf("FindContent ran an FCS query but it should have run the fallback")
	}
	if ops.ContentIdQueried != 2222 {
		t.Errorf("FindContent should have queried content ID 2222 but we got %d", ops.ContentIdQueried)
	}

	if ops.IdMappingIndexQueried != IdMappingIndexOctid {
		t.Errorf("FindContent should have queried id mappings by filebase but we got %s", ops.IdMappingIndexQueried)
	}
	if ops.IdMappingKeyFieldQueried != IdMappingKeyfieldOctid {
		t.Errorf("FindContent should have queried id mappings on field filebase but we got %s", ops.IdMappingKeyFieldQueried)
	}
	if ops.IdMappingSearchTermQueried != int64(123456) {
		t.Errorf("FindContent queried ID mappings on the wrong term, got %d", ops.IdMappingSearchTermQueried)
	}
}
