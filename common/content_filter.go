package common

/*
RemoveIndex Output an array of pointers to Encoding after removing the pointer to Encoding at the supplied index
Arguments:
- encodings - An array of pointers to Encoding
- index - The index to remove
Returns:
- Array of pointers to Encoding
*/
func RemoveIndex(encodings []*Encoding, index int) []*Encoding {
	return append(encodings[:index], encodings[index+1:]...)
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
	for index, element := range encodings {
		if element.Format != format {
			encodings = RemoveIndex(encodings, index)
			continue
		}

		if element.Mobile != need_mobile {
			encodings = RemoveIndex(encodings, index)
			continue
		}

		if element.VBitrate < minbitrate {
			encodings = RemoveIndex(encodings, index)
			continue
		}

		if element.VBitrate > maxbitrate {
			encodings = RemoveIndex(encodings, index)
			continue
		}

		if element.FrameHeight < minheight {
			encodings = RemoveIndex(encodings, index)
			continue
		}

		if element.FrameHeight > maxheight {
			encodings = RemoveIndex(encodings, index)
			continue
		}

		if element.FrameWidth < minwidth {
			encodings = RemoveIndex(encodings, index)
			continue
		}

		if element.FrameWidth > maxwidth {
			encodings = RemoveIndex(encodings, index)
			continue
		}
	}

	return encodings
}
