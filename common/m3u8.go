package common

import (
	"log"
	"regexp"
)

/*
HasDodgyM3U8Format fixes a problem with iOS devices supplying malformed format strings to the endpoint.
There is a bug in iOS clients whereby rather than using the _actual_ redirected m3u8 URL to locate submanifests
it simply takes the orignal referer and replaces everything after the last /.  So, if the client came here through
us the url becomes endpoint.yadayada.com/interactivevideos/video.php?format=video/{filename}.m3u8.
We deem these "dodgy m3u8 format strings" and deal with them here by supplying override values back to the main func.

Returns:

- a string with the correct filename value to use

- a boolean indicating whether the fix should be applied or not. If `true` then this is a "dodgy" format string that must be fixed
*/
func HasDodgyM3U8Format(format string) (string, bool) {
	matcher := regexp.MustCompile("video/(.*\\.m3u8)$")
	results := matcher.FindAllStringSubmatch(format, -1)
	if results != nil {
		log.Printf("DEBUG HasDodgyM3U8Format got matches %v", results)
		return results[0][1], true
	} else {
		return "", false
	}
}
