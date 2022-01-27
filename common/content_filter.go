package common

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
func TestEncoding(encoding *Encoding, format string, need_mobile bool, minbitrate int32, maxbitrate int32, minheight int32, maxheight int32, minwidth int32, maxwidth int32) bool {
	if encoding.Format != format {
		return false
	}

	if encoding.Mobile != need_mobile {
		return false
	}

	if encoding.VBitrate < minbitrate {
		return false
	}

	if encoding.VBitrate > maxbitrate {
		return false
	}

	if encoding.FrameHeight < minheight {
		return false
	}

	if encoding.FrameHeight > maxheight {
		return false
	}

	if encoding.FrameWidth < minwidth {
		return false
	}

	if encoding.FrameWidth > maxwidth {
		return false
	}

	return true
}

/*
ContentFilter Output a filtered array of pointers to Encoding based on the other arguments
Arguments:
- encodings - An array of pointers to Encoding
- format - The required format
- need_mobile - Set this to true if a mobile encoding is required
- minbitrate - The minimum required bit rate
- maxbitrate - The maximum required bit rate
- minheight - The minimum required frame height
- maxheight - The maximum required frame height
- minwidth - The minimum required frame width
- maxwidth - The maximum required frame width
Returns:
- Array of pointers to Encoding
*/
func ContentFilter(encodings []*Encoding, format string, need_mobile bool, minbitrate int32, maxbitrate int32, minheight int32, maxheight int32, minwidth int32, maxwidth int32) []*Encoding {
	var encodingsToReturn []*Encoding
	for _, element := range encodings {
		if TestEncoding(element, format, need_mobile, minbitrate, maxbitrate, minheight, maxheight, minwidth, maxwidth) {
			encodingsToReturn = append(encodingsToReturn, element)
		}
	}

	return encodingsToReturn
}
