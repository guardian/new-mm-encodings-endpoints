package common

import (
	"errors"
	"reflect"
)

type PosterFrame struct {
	PosterId   int32  `json:"poster_id"`
	EncodingId int32  `json:"encoding_id"`
	ContentId  int32  `json:"content_id"`
	PosterUrl  string `json:"poster_url"`
	MimeType   string `json:"mime_type"`
}

func PosterFrameFromDynamo(rec *RawDynamoRecord) (result *PosterFrame, e error) {
	defer func() {
		//we allow the routine to panic if the typecast below fails then recover it here and return an error.
		//the underlying cause should have been logged out already.
		if r := recover(); r != nil {
			result = nil
			e = errors.New("the given record is not a ContentResult")
		}
	}()

	newRecord := &PosterFrame{
		PosterId:   extractDynamoField(rec, "posterid", reflect.Int32, false).(int32),
		EncodingId: extractDynamoField(rec, "encodingid", reflect.Int32, false).(int32),
		ContentId:  extractDynamoField(rec, "contentid", reflect.Int32, false).(int32),
		PosterUrl:  extractDynamoField(rec, "poster_url", reflect.String, false).(string),
		MimeType:   extractDynamoField(rec, "mime_type", reflect.String, false).(string),
	}
	return newRecord, nil
}
