package common

import (
	"regexp"
)

/*
ForceHTTPS Output a URL starting with https unless allowInsecure is set to true
Arguments:
- url - URL to process
- allowInsecure - If set to true, causes the function to output the input URL
Returns:
- URL
*/
func ForceHTTPS(url string, allowInsecure bool) string {
	if allowInsecure {
		return url
	} else {
		var re = regexp.MustCompile(`^htt?p:`)
		urlToReturn := re.ReplaceAllString(url, `https:`)
		return urlToReturn
	}
}
