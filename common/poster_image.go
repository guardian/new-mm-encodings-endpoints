package common

import (
	"errors"
	"regexp"
)

/*
GeneratePosterImageURL Outputs the correct URL for the poster image
Arguments:
- url - The URL of the media file
Returns:
- URL of the poster image
*/
func GeneratePosterImageURL(url string) (string, error) {
	var re = regexp.MustCompile(`^(.*)\.[^\.]+$`)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return "", errors.New("the CDN URL was malformed")
	}
	return matches[1] + "_poster.jpg", nil
}
