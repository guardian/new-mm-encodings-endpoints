package common

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
	"strconv"
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
extractContentResultField is an internal method that grabs the relevant field from Dynamodb, casts it to the right type
and returns it as an interface.

The return value is `nil` if either the field does not exist or is the wrong type.
If a value is returned it should always be castable to the type given in the `t` parameter without further checks
*/
func extractContentResultField(rec *RawDynamoRecord, fieldName string, t reflect.Kind, nullable bool) interface{} {
	if dynamoValue, valueExists := (*rec)[fieldName]; valueExists {
		switch t {
		case reflect.String:
			if v, isRightType := dynamoValue.(*types.AttributeValueMemberS); isRightType {
				return v.Value
			} else {
				log.Printf("WARNING Field %s was not castable to a string", fieldName)
				return nil
			}
		case reflect.Int32:
			if v, isRightType := dynamoValue.(*types.AttributeValueMemberN); isRightType {
				intval, err := strconv.ParseInt(v.Value, 10, 32)
				if err != nil {
					log.Printf("WARNING Field %s value %s could not be converted to int32: %s", fieldName, v.Value, err)
					return nil
				}
				return int32(intval)
			} else {
				log.Printf("WARNING Field %s was not castable to a number", fieldName)
				return nil
			}
		case reflect.Float32:
			if v, isRightType := dynamoValue.(*types.AttributeValueMemberN); isRightType {
				intval, err := strconv.ParseFloat(v.Value, 32)
				if err != nil {
					log.Printf("WARNING Field %s value %s could not be converted to float32: %s", fieldName, v.Value, err)
					return nil
				}
				return float32(intval)
			} else {
				log.Printf("WARNING Field %s was not castable to a number", fieldName)
				return nil
			}
		case reflect.Bool:
			if v, isRightType := dynamoValue.(*types.AttributeValueMemberBOOL); isRightType {
				return v.Value
			} else {
				log.Printf("WARNING Field %s was not castable to a bool", fieldName)
				return nil
			}
		default:
			log.Printf("WARNING Field %s has a type of %s which is not handled", fieldName, reflect.TypeOf(dynamoValue))
			return nil
		}
	} else {
		if nullable {
			switch t {
			case reflect.String:
				return ""
			case reflect.Int32:
				return int32(0)
			case reflect.Int64:
				return int64(0)
			case reflect.Float32:
				return float32(0)
			case reflect.Bool:
				return false
			default:
				return nil
			}
		} else {
			log.Printf("WARNING Field %s does not exist on the incoming record", fieldName)
			return nil
		}
	}
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

	lastUpdateTime, err := time.Parse(time.RFC3339, extractContentResultField(rec, "lastupdate", reflect.String, false).(string))
	if err != nil {
		log.Printf("WARNING Field 'lastupdate' is not a valid timestamp: %s", err)
		return nil, errors.New("invalid timestamp")
	}

	newRecord := &Encoding{
		EncodingId:  extractContentResultField(rec, "encodingid", reflect.Int32, false).(int32),
		ContentId:   extractContentResultField(rec, "contentid", reflect.Int32, false).(int32),
		Url:         extractContentResultField(rec, "url", reflect.String, false).(string),
		Format:      extractContentResultField(rec, "format", reflect.String, false).(string),
		Mobile:      extractContentResultField(rec, "mobile", reflect.Bool, false).(bool),
		Multirate:   extractContentResultField(rec, "multirate", reflect.Bool, false).(bool),
		VCodec:      extractContentResultField(rec, "vcodec", reflect.String, true).(string),
		ACodec:      extractContentResultField(rec, "acodec", reflect.String, true).(string),
		VBitrate:    extractContentResultField(rec, "vbitrate", reflect.Int32, true).(int32),
		ABitrate:    extractContentResultField(rec, "abitrate", reflect.Int32, true).(int32),
		LastUpdate:  lastUpdateTime,
		FrameWidth:  extractContentResultField(rec, "frame_width", reflect.Int32, false).(int32),
		FrameHeight: extractContentResultField(rec, "frame_height", reflect.Int32, false).(int32),
		Duration:    extractContentResultField(rec, "duration", reflect.Float32, false).(float32),
		FileSize:    extractContentResultField(rec, "file_size", reflect.Int64, false).(int64),
		FCSID:       extractContentResultField(rec, "fcs_id", reflect.String, true).(string),
		OctopusId:   extractContentResultField(rec, "octopus_id", reflect.Int32, true).(int32),
		Aspect:      extractContentResultField(rec, "aspect", reflect.String, false).(string),
	}

	return newRecord, nil
}
