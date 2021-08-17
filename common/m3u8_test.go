package common

import "testing"

func TestHasDodgyM3U8FormatRequired(t *testing.T) {
	result := HasDodgyM3U8Format("video/somefilename.m3u8")
	if result == nil {
		t.Error("HasDodgyM3U8Format returned nil when it should have given data")
	} else {
		if result["format"] != "video/m3u8" {
			t.Errorf("HasDodgyM3U8Format returned wrong format %s, expected video/m3u8", result["format"])
		}
		if result["filename"] != "somefilename.m3u8" {
			t.Errorf("HasDodgyM3U8Format returned wrong filename %s, expected somefilename.m3u8", result["filename"])
		}

	}
}

func TestHasDodgyM3U8FormatNotRequired(t *testing.T) {
	result := HasDodgyM3U8Format("video/mp4")
	if result != nil {
		t.Errorf("HasDodgyM3U8Format returned %v when it should have given data", result)
	}
}
