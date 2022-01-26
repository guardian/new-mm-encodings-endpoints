package common

func RemoveIndex(s []*Encoding, index int) []*Encoding {
	return append(s[:index], s[index+1:]...)
}

func ContentFilter(encodings []*Encoding, format string, need_mobile bool, minbitrate int32, maxbitrate int32, minheight int32, maxheight int32, minwidth int32, maxwidth int32) (encodings_to_return []*Encoding) {
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
