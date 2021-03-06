package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"reflect"
	"testing"
	"time"
)

func ReturnCorrectTimeObject() time.Time {
	timeObject, err := time.Parse(time.RFC3339, "2016-02-03T04:05:06Z")
	if err != nil {
		fmt.Printf("WARNING Could not create time object: %s", err)
	}
	return timeObject
}

/*
Tests that the code which removes records with the wrong format works
*/
func TestContentFilterFormat(t *testing.T) {
	formats := []string{"test"}
	testarray := []*Encoding{{1, 1, "test", "test2", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, false, 1, 1, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

/*
Should return data that has _any_ of the given formats specified
*/
func TestContentFilterAlternateFormat(t *testing.T) {
	formats := []string{"socks", "test"}
	testarray := []*Encoding{{1, 1, "test", "test2", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, false, 1, 1, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with the wrong mobile setting works
*/
func TestContentFilterMobile(t *testing.T) {
	formats := []string{"test"}
	testarray := []*Encoding{{1, 1, "test", "test", true, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", true, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, true, 1, 1, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too low a bit rate works
*/
func TestContentFilterMinBitRate(t *testing.T) {
	formats := []string{"test"}
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, false, 3000, 6000, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too high a bit rate works
*/
func TestContentFilterMaxBitRate(t *testing.T) {
	formats := []string{"test"}
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 8000, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, false, 3000, 6000, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too low a frame height works
*/
func TestContentFilterMinHeight(t *testing.T) {
	formats := []string{"test"}
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1000, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 4000, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1000, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, false, 3000, 6000, 800, 2000, 1, 1)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too high a frame height works
*/
func TestContentFilterMaxHeight(t *testing.T) {
	formats := []string{"test"}
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1000, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 4000, 1, ReturnCorrectTimeObject(), 1, 4000, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1000, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, false, 3000, 6000, 800, 2000, 1, 1)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too low a frame width works
*/
func TestContentFilterMinWidth(t *testing.T) {
	formats := []string{"test"}
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1000, 1000, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 4000, 1, ReturnCorrectTimeObject(), 800, 1000, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1000, 1000, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, false, 3000, 6000, 800, 2000, 900, 2000)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too high a frame width works
*/
func TestContentFilterMaxWidth(t *testing.T) {
	formats := []string{"test"}
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1000, 1000, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 4000, 1, ReturnCorrectTimeObject(), 3000, 1000, 4.0, 1, "test", 1, "test"}}
	expectedOutput := &ContentResult{Encoding{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1000, 1000, 4.0, 1, "test", 1, "test"}, "", ""}
	result := ContentFilter(testarray, &formats, false, 3000, 6000, 800, 2000, 900, 2000)
	if !reflect.DeepEqual(result, expectedOutput) {
		t.Errorf("Unexpected output")
	}
}

func TestIsStringInList(t *testing.T) {
	if isStringInList(aws.String("audio/mpeg"), &[]string{"audio/mpeg", "audio/mp3"}) == false {
		t.Error("isStringInList failed to detect 'audio/mpeg' in a string")
	}
}
