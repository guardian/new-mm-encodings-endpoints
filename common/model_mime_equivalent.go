package common

import (
	"errors"
	"reflect"
)

type MimeEquivalent struct {
	Id             int32
	RealName       string
	MimeEquivalent string
}

func MimeEquivalentFromDynamo(rec *RawDynamoRecord) (result *MimeEquivalent, e error) {
	defer func() {
		//we allow the routine to panic if the typecast below fails then recover it here and return an error.
		//the underlying cause should have been logged out already.
		if r := recover(); r != nil {
			result = nil
			e = errors.New("the given record is not a ContentResult")
		}
	}()

	newRecord := &MimeEquivalent{
		Id:             extractDynamoField(rec, "id", reflect.Int32, true).(int32),
		RealName:       extractDynamoField(rec, "real_name", reflect.String, false).(string),
		MimeEquivalent: extractDynamoField(rec, "mime_equivalent", reflect.String, false).(string),
	}
	return newRecord, nil
}
