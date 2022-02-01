package common

/*
isStringInList will return `true` if the given string is in the given list or `false` otherwise
*/
func isStringInList(needle *string, haystack *[]string) bool {
	for _, s := range *haystack {
		if s == *needle {
			return true
		}
	}
	return false
}

/*
TestEncoding Output true if the encoding should pass the filter and false if it should not
Arguments:
- encoding - A pointer to Encoding
- format - The required format
- need_mobile - Set this to true if a mobile encoding is required
- minbitrate - The minimum required bit rate
- maxbitrate - The maximum required bit rate
- minheight - The minimum required frame height
- maxheight - The maximum required frame height
- minwidth - The minimum required frame width
- maxwidth - The maximum required frame width
Returns:
- bool - true if the encoding should pass and false if it should not
*/
func TestEncoding(encoding *Encoding, formats *[]string, need_mobile bool, minbitrate int32, maxbitrate int32, minheight int32, maxheight int32, minwidth int32, maxwidth int32) bool {
	if len(*formats) > 0 && !isStringInList(&encoding.Format, formats) {
		return false
	}

	if encoding.Mobile != need_mobile {
		return false
	}

	if (encoding.VBitrate < minbitrate) && (minbitrate != 0) {
		return false
	}

	if (encoding.VBitrate > maxbitrate) && (maxbitrate != 0) {
		return false
	}

	if (encoding.FrameHeight < minheight) && (minheight != 0) {
		return false
	}

	if (encoding.FrameHeight > maxheight) && (maxheight != 0) {
		return false
	}

	if (encoding.FrameWidth < minwidth) && (minwidth != 0) {
		return false
	}

	if (encoding.FrameWidth > maxwidth) && (maxwidth != 0) {
		return false
	}

	return true
}

/*
ContentFilter Output a pointer to a ContentResult object after filtering an array of pointers to Encoding based on the other arguments
Arguments:
- encodings - An array of pointers to Encoding
- formats - An array of 0 or more strings representing a MIME type. If this is non-zero-length then at least one of the format strings must match.
- need_mobile - Set this to true if a mobile encoding is required
- minbitrate - The minimum required bit rate
- maxbitrate - The maximum required bit rate
- minheight - The minimum required frame height
- maxheight - The maximum required frame height
- minwidth - The minimum required frame width
- maxwidth - The maximum required frame width
Returns:
- ContentResult object populated with the best pointer to an Encoding
*/
func ContentFilter(encodings []*Encoding, formats *[]string, need_mobile bool, minbitrate int32, maxbitrate int32, minheight int32, maxheight int32, minwidth int32, maxwidth int32) *ContentResult {
	var encodingsToReturn []*Encoding
	for _, element := range encodings {
		if TestEncoding(element, formats, need_mobile, minbitrate, maxbitrate, minheight, maxheight, minwidth, maxwidth) {
			encodingsToReturn = append(encodingsToReturn, element)
		}
	}

	if len(encodingsToReturn) == 0 {
		return nil
	} else {
		return &ContentResult{*encodingsToReturn[0], "", ""}
	}
}
