package main

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/guardian/new-encodings-endpoints/common"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type EndpointEvent struct {
	Uid                   uuid.UUID
	Timestamp             *time.Time
	AccessUrl             string
	ExpectedOutputMessage string
	ExpectedOutputHeaders map[string]string
	ExpectedResponse      int16
}

/*
IsValid will return `true` if the EndpointEvent has data that looks right
*/
func (e *EndpointEvent) IsValid() bool {
	emptyUid := uuid.UUID{}
	return e.Uid != emptyUid && e.Timestamp != nil && e.AccessUrl != "" && e.ExpectedResponse >= 200
}

func (e *EndpointEvent) FormattedTimestamp() string {
	return e.Timestamp.Format(time.RFC3339)
}

func getStringValue(d types.AttributeValue) (string, error) {
	if stringVal, haveStringVal := d.(*types.AttributeValueMemberS); haveStringVal {
		return stringVal.Value, nil
	} else {
		return "", errors.New("value was not a string")
	}
}

func getUidValue(d types.AttributeValue) (uuid.UUID, error) {
	stringVal, err := getStringValue(d)
	if err != nil {
		return uuid.UUID{}, err
	}
	uid, err := uuid.Parse(stringVal)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uid, nil
}

func getTimeValue(d types.AttributeValue) (*time.Time, error) {
	stringVal, err := getStringValue(d)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, stringVal)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

var HeaderSplit = regexp.MustCompile("^([^:]+):\\s*(.*)$")

func EndpointEventFromDynamo(rec *common.RawDynamoRecord) (*EndpointEvent, error) {
	evt := &EndpointEvent{}

	if uidString, haveUid := (*rec)["uid"]; haveUid {
		uid, err := getUidValue(uidString)
		if err != nil {
			log.Printf("ERROR Invalid EndpointEvent ID: %s", err)
			return nil, err
		}
		evt.Uid = uid
	}
	if timestamp, haveTimestamp := (*rec)["timestamp"]; haveTimestamp {
		t, err := getTimeValue(timestamp)
		if err != nil {
			log.Printf("ERROR Invalid EndpointEvent timestamp: %s", err)
			return nil, err
		}
		evt.Timestamp = t
	}
	if url, haveUrl := (*rec)["access_url"]; haveUrl {
		u, err := getStringValue(url)
		if err != nil {
			log.Printf("ERROR Invalid EndpointEvent URL: %s", err)
			return nil, err
		}
		evt.AccessUrl = u
	}
	if msg, haveMsg := (*rec)["output_message"]; haveMsg {
		if _, isNull := msg.(*types.AttributeValueMemberNULL); !isNull { //null values allowed here
			m, err := getStringValue(msg)
			if err != nil {
				log.Printf("ERROR Invalid EndpointEvent message: %s", err)
			}
			evt.ExpectedOutputMessage = m
		}
	}
	if h, haveHeaders := (*rec)["php_headers"]; haveHeaders {
		vals, isList := h.(*types.AttributeValueMemberL)
		if isList {
			evt.ExpectedOutputHeaders = make(map[string]string, len(vals.Value))
			n := 0
			for _, v := range vals.Value {
				if s, isString := v.(*types.AttributeValueMemberS); isString {
					matches := HeaderSplit.FindAllStringSubmatch(s.Value, -1)
					if matches != nil {
						n++
						evt.ExpectedOutputHeaders[matches[0][1]] = matches[0][2]
					}
				}
			}
		} else {
			log.Printf("ERROR Invalid EndpointEvent: `php_headers` is not a List, got %s", reflect.TypeOf(h))
		}
	}
	if r, haveResponseCode := (*rec)["response_code"]; haveResponseCode {
		val, isNumber := r.(*types.AttributeValueMemberN)
		if isNumber {
			intVal, _ := strconv.ParseInt(val.Value, 10, 16)
			evt.ExpectedResponse = int16(intVal)
		} else {
			log.Printf("ERROR Invalid EndpointEvent: `response_code` is not a Number")
		}
	}
	if evt.IsValid() {
		return evt, nil
	} else {
		return nil, errors.New("EndpointEvent record is not valid")
	}
}
