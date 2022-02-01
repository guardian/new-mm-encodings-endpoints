package common

import (
	"reflect"
	"testing"
)

/*
Tests MakeResponseRedirect outputs the correct data
*/
func TestMakeResponseRedirect(t *testing.T) {
	result := MakeResponseRedirect("https://test.url/")
	expectedOutput := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET, OPTIONS",
		"Access-Control-Allow-Headers":     "*",
		"Access-Control-Allow-Credentials": "false",
		"Access-Control-Max-Age":           "3600",
		"Location":                         "https://test.url/",
	}
	if !reflect.DeepEqual(result.Headers, expectedOutput) {
		t.Errorf("Unexpected output: %s", result.Headers)
	}
	if !reflect.DeepEqual(result.StatusCode, 302) {
		t.Errorf("Unexpected output: %d", result.StatusCode)
	}
	if !reflect.DeepEqual(result.Body, "") {
		t.Errorf("Unexpected output: %s", result.Body)
	}
}
