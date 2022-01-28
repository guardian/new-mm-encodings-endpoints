package common

import "testing"

func TestHasDodgyM3U8FormatRequired(t *testing.T) {
	updatedFilename, isDodgy := HasDodgyM3U8Format("video/somefilename.m3u8")
	if !isDodgy {
		t.Error("HasDodgyM3U8Format failed to pick up a dodgy format string")
	}

	if updatedFilename != "somefilename.m3u8" {
		t.Errorf("HasDodgyM3U8Format returned incorrect filename override: %s", updatedFilename)
	}
}

func TestHasDodgyM3U8FormatNotRequired(t *testing.T) {
	_, isDodgy := HasDodgyM3U8Format("video/mp4")
	if isDodgy {
		t.Errorf("HasDodgyM3U8Format counted a valid format string as dodgy")
	}
}
