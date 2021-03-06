package common

import (
	"testing"
)

/*
Tests that the correct URL is returned for the poster image
*/
func TestGeneratePosterImageURL(t *testing.T) {
	result, _ := GeneratePosterImageURL("https://cdn.theguardian.tv/HLS/2018/06/06/091101BangladeshVillages.m3u8", false)
	if result != "https://cdn.theguardian.tv/HLS/2018/06/06/091101BangladeshVillages_poster.jpg" {
		t.Errorf("Unexpected output: %s", result)
	}
	result, _ = GeneratePosterImageURL("https://cdn.theguardian.tv/HLS/2018/06/06/091101BangladeshVillages.m3u8", true)
	if result != "https://cdn.theguardian.tv/HLS/2018/06/06/091101BangladeshVillages_poster.png" {
		t.Errorf("Unexpected output: %s", result)
	}
}

/*
Tests an error is returned if the URL can not be parsed
*/
func TestGeneratePosterImageURLError(t *testing.T) {
	_, err := GeneratePosterImageURL("test", false)
	if err == nil {
		t.Error("GeneratePosterImageURL returned no error for an invalid URL")
	}
}
