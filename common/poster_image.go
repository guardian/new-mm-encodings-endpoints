package common

import (
	"regexp"
)

/*
GeneratePosterImageURL Outputs the correct URL for the poster image
Arguments:
- url - The URL of the media file
Returns:
- URL of the poster image
*/
func GeneratePosterImageURL(url string) string {
	var re = regexp.MustCompile(`^(.*)\.[^\.]+$`)
	matches := re.FindStringSubmatch(url)
	return matches[1] + "_poster.jpg"
}
