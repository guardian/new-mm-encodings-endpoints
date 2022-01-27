package common

import "testing"

/*
Tests that the correct URL is returned for the poster image
*/
func TestGeneratePosterImageURL(t *testing.T) {
	result := GeneratePosterImageURL("https://cdn.theguardian.tv/HLS/2018/06/06/091101BangladeshVillages.m3u8")
	if result != "https://cdn.theguardian.tv/HLS/2018/06/06/091101BangladeshVillages_poster.jpg" {
		t.Errorf("Unexpected output: %s", result)
	}
}
