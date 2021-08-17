package common

import (
	"log"
	"regexp"
)

/**
HasDodgyM3U8 format fixes a problem with iOS devices supplying malformed format strings to the endpoint.
There is a bug in iOS clients whereby rather than using the _actual_ redirected m3u8 URL to locate submanifests
it simply takes the orignal referer and replaces everything after the last /.  So, if the client came here through
us the url becomes endpoint.yadayada.com/interactivevideos/video.php?format=video/{filename}.m3u8.
We deem these "dodgy m3u8 format strings" and deal with them here by supplying override values back to the main func.

Returns a map of data to override the existing data if required, otherwise nil.
*/
func HasDodgyM3U8Format(format string) map[string]string {
	matcher := regexp.MustCompile("video/(.*\\.m3u8)$")
	results := matcher.FindAllStringSubmatch(format, -1)
	if results != nil {
		log.Printf("DEBUG HasDodgyM3U8Format got matches %v", results)
		return map[string]string{
			"format":   "video/m3u8",
			"filename": results[0][1],
		}
	} else {
		return nil
	}
}
