package common

import (
	"errors"
	"log"
	"reflect"
	"time"
)

/*
Encoding represents a row from the `Encodings` table
*/
type Encoding struct {
	EncodingId  int32     `json:"encoding_id"` //NOT NULL
	ContentId   int32     `json:"content_id"`  //NOT NULL
	Url         string    `json:"url"`         //NOT NULL
	Format      string    `json:"format"`      //NOT NULL
	Mobile      bool      `json:"mobile"`      //NOT NULL
	Multirate   bool      `json:"multirate"`   //NOT NULL
	VCodec      string    `json:"vcodec"`
	ACodec      string    `json:"acodec"`
	VBitrate    int32     `json:"vbitrate"`
	ABitrate    int32     `json:"abitrate"`
	LastUpdate  time.Time `json:"last_update"`  //NOT NULL, defaults to current time
	FrameWidth  int32     `json:"frame_width"`  //NOT NULL
	FrameHeight int32     `json:"frame_height"` //NOT NULL
	Duration    float32   `json:"duration"`     //NOT NULL
	FileSize    int64     `json:"file_size"`    //NOT NULL
	FCSID       string    `json:"fcs_id"`       //NOT NULL
	OctopusId   int32     `json:"octopus_id"`   //NOT NULL aka 'title id'
	Aspect      string    `json:"aspect"`       //NOT NULL
}

/*
EncodingFromDynamo takes a RawDynamoRecord (aka map of string -> dynamo value) and marshals it into a EncodingFromDynamo
structure.  If we can't validate the structure an error is returned instead.
Arguments: - rec - pointer to a raw dynamo record
Returns:
- pointer to the created ContentResult or nil on error
- nil or an `error` if an error occurs
*/
func EncodingFromDynamo(rec *RawDynamoRecord) (result *Encoding, e error) {
	defer func() {
		//we allow the routine to panic if the typecast below fails then recover it here and return an error.
		//the underlying cause should have been logged out already.
		if r := recover(); r != nil {
			result = nil
			e = errors.New("the given record is not a ContentResult")
		}
	}()

	lastUpdateTime, err := time.Parse(time.RFC3339, extractDynamoField(rec, "lastupdate", reflect.String, false).(string))
	if err != nil {
		log.Printf("WARNING Field 'lastupdate' is not a valid timestamp: %s", err)
		return nil, errors.New("invalid timestamp")
	}

	newRecord := &Encoding{
		EncodingId:  extractDynamoField(rec, "encodingid", reflect.Int32, false).(int32),
		ContentId:   extractDynamoField(rec, "contentid", reflect.Int32, false).(int32),
		Url:         extractDynamoField(rec, "url", reflect.String, false).(string),
		Format:      extractDynamoField(rec, "format", reflect.String, false).(string),
		Mobile:      extractDynamoField(rec, "mobile", reflect.Bool, false).(bool),
		Multirate:   extractDynamoField(rec, "multirate", reflect.Bool, false).(bool),
		VCodec:      extractDynamoField(rec, "vcodec", reflect.String, true).(string),
		ACodec:      extractDynamoField(rec, "acodec", reflect.String, true).(string),
		VBitrate:    extractDynamoField(rec, "vbitrate", reflect.Int32, true).(int32),
		ABitrate:    extractDynamoField(rec, "abitrate", reflect.Int32, true).(int32),
		LastUpdate:  lastUpdateTime,
		FrameWidth:  extractDynamoField(rec, "frame_width", reflect.Int32, false).(int32),
		FrameHeight: extractDynamoField(rec, "frame_height", reflect.Int32, false).(int32),
		Duration:    extractDynamoField(rec, "duration", reflect.Float32, false).(float32),
		FileSize:    extractDynamoField(rec, "file_size", reflect.Int64, false).(int64),
		FCSID:       extractDynamoField(rec, "fcs_id", reflect.String, true).(string),
		OctopusId:   extractDynamoField(rec, "octopus_id", reflect.Int32, true).(int32),
		Aspect:      extractDynamoField(rec, "aspect", reflect.String, false).(string),
	}

	return newRecord, nil
}
