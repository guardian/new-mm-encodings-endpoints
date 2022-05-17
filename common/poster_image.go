package common

import (
	"errors"
	"log"
	"regexp"
)

/*
GeneratePosterImageURL outputs the likely URL for a poster image corresponding to this file. Note that the validity of the URL is not checked so it may return a 404 if no poster image was created in the first place.
Arguments:
- url - The URL of the media file
Returns:
- URL of the poster image
*/
func GeneratePosterImageURL(url string, pngPoster bool) (string, error) {
	var re = regexp.MustCompile(`^(.*)\.[^\.]+$`)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return "", errors.New("the CDN URL was malformed (no file extension) ")
	}

	log.Printf("GeneratePosterImageURL: pngPoster is %t", pngPoster)
	xtn := ".jpg"
	if pngPoster {
		xtn = ".png"
	}
	log.Printf("poster path is %s", matches[1]+"_poster"+xtn)

	return matches[1] + "_poster" + xtn, nil
}
