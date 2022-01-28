package common

import (
	"errors"
	"regexp"
)

/*
GeneratePosterImageURL outputs the likely URL for a poster image corresponding to this file. Note that the validity of the URL is not checked so it may return a 404 if no poster image was created in the first place.
Arguments:
- url - The URL of the media file
Returns:
- URL of the poster image
*/
func GeneratePosterImageURL(url string) (string, error) {
	var re = regexp.MustCompile(`^(.*)\.[^\.]+$`)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return "", errors.New("the CDN URL was malformed (no file extension) ")
	}
	return matches[1] + "_poster.jpg", nil
}
