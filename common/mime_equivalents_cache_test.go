package common

import (
	"reflect"
	"testing"
)

func TestMimeEquivalentsCacheImpl_EquivalentsFor(t *testing.T) {
	toTest := &MimeEquivalentsCacheImpl{
		loadedData: map[string]string{
			"video/m3u8":            "application/x-mpegURL",
			"application/x-mpegURL": "video/m3u8",
		},
	}

	m3u8result := toTest.EquivalentsFor("video/m3u8")
	if !reflect.DeepEqual(m3u8result, []string{"video/m3u8", "application/x-mpegURL"}) {
		t.Errorf("EquivalentsFor returned %v when we expected \"video/m3u8\", \"application/x-mpegURL\"", m3u8result)
	}

	mpurlResult := toTest.EquivalentsFor("application/x-mpegURL")
	if !reflect.DeepEqual(mpurlResult, []string{"application/x-mpegURL", "video/m3u8"}) {
		t.Errorf("EquivalentsFor returned %v when we expected \"application/x-mpegURL\", \"video/m3u8\"", mpurlResult)
	}

	otherResult := toTest.EquivalentsFor("media/rhubarb")
	if !reflect.DeepEqual(otherResult, []string{"media/rhubarb"}) {
		t.Errorf("EquivalentsFor returned %v when we expected \"media/rhubarb\"", otherResult)
	}
}
