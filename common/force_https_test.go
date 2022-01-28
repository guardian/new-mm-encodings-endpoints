package common

import "testing"

/*
Tests that the correct URL is returned when http is input
*/
func TestForceHTTPS(t *testing.T) {
	result := ForceHTTPS("http://test/", false)
	if result != "https://test/" {
		t.Errorf("Unexpected output: %s", result)
	}
}

/*
Tests that the correct URL is returned when https is input
*/
func TestForceHTTPSCanCopeWithHTTPS(t *testing.T) {
	result := ForceHTTPS("https://test/", false)
	if result != "https://test/" {
		t.Errorf("Unexpected output: %s", result)
	}
}

/*
Tests that the correct URL is returned when allowInsecure is true
*/
func TestForceHTTPSCanAllowInsecure(t *testing.T) {
	result := ForceHTTPS("http://test/", true)
	if result != "http://test/" {
		t.Errorf("Unexpected output: %s", result)
	}
}
