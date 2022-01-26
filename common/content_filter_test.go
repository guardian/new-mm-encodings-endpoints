package common

import (
	"fmt"
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
	testarray := []*Encoding{{1, 1, "test", "test2", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	result := ContentFilter(testarray, "test", false, 1, 1, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedarray) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with the wrong mobile setting works
*/
func TestContentFilterMobile(t *testing.T) {
	testarray := []*Encoding{{1, 1, "test", "test", true, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedarray := []*Encoding{{1, 1, "test", "test", true, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	result := ContentFilter(testarray, "test", true, 1, 1, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedarray) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too low a bit rate works
*/
func TestContentFilterMinBitRate(t *testing.T) {
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 1, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	result := ContentFilter(testarray, "test", false, 3000, 6000, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedarray) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too high a bit rate works
*/
func TestContentFilterMaxBitRate(t *testing.T) {
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 8000, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	result := ContentFilter(testarray, "test", false, 3000, 6000, 1, 1, 1, 1)
	if !reflect.DeepEqual(result, expectedarray) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too low a frame height works
*/
func TestContentFilterMinHeight(t *testing.T) {
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1000, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 4000, 1, ReturnCorrectTimeObject(), 1, 1, 4.0, 1, "test", 1, "test"}}
	expectedarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1000, 4.0, 1, "test", 1, "test"}}
	result := ContentFilter(testarray, "test", false, 3000, 6000, 800, 2000, 1, 1)
	if !reflect.DeepEqual(result, expectedarray) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too high a frame height works
*/
func TestContentFilterMaxHeight(t *testing.T) {
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1000, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 4000, 1, ReturnCorrectTimeObject(), 1, 4000, 4.0, 1, "test", 1, "test"}}
	expectedarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1, 1000, 4.0, 1, "test", 1, "test"}}
	result := ContentFilter(testarray, "test", false, 3000, 6000, 800, 2000, 1, 1)
	if !reflect.DeepEqual(result, expectedarray) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too low a frame width works
*/
func TestContentFilterMinWidth(t *testing.T) {
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1000, 1000, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 4000, 1, ReturnCorrectTimeObject(), 800, 1000, 4.0, 1, "test", 1, "test"}}
	expectedarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1000, 1000, 4.0, 1, "test", 1, "test"}}
	result := ContentFilter(testarray, "test", false, 3000, 6000, 800, 2000, 900, 2000)
	if !reflect.DeepEqual(result, expectedarray) {
		t.Errorf("Unexpected output")
	}
}

/*
Tests that the code which removes records with too high a frame width works
*/
func TestContentFilterMaxWidth(t *testing.T) {
	testarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1000, 1000, 4.0, 1, "test", 1, "test"}, {1, 1, "test", "test", false, false, "test", "test", 4000, 1, ReturnCorrectTimeObject(), 3000, 1000, 4.0, 1, "test", 1, "test"}}
	expectedarray := []*Encoding{{1, 1, "test", "test", false, false, "test", "test", 3557, 1, ReturnCorrectTimeObject(), 1000, 1000, 4.0, 1, "test", 1, "test"}}
	result := ContentFilter(testarray, "test", false, 3000, 6000, 800, 2000, 900, 2000)
	if !reflect.DeepEqual(result, expectedarray) {
		t.Errorf("Unexpected output")
	}
}
