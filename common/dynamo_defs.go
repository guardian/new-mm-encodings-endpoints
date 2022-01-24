package common

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
	"strconv"
)

type RawDynamoRecord map[string]types.AttributeValue

/*
extractDynamoField is an internal method that grabs the relevant field from Dynamodb, casts it to the right type
and returns it as an interface.

The return value is `nil` if either the field does not exist or is the wrong type.
If a value is returned it should always be castable to the type given in the `t` parameter without further checks
*/
func extractDynamoField(rec *RawDynamoRecord, fieldName string, t reflect.Kind, nullable bool) interface{} {
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
		case reflect.Int64:
			if v, isRightType := dynamoValue.(*types.AttributeValueMemberN); isRightType {
				intval, err := strconv.ParseInt(v.Value, 10, 64)
				if err != nil {
					log.Printf("WARNING Field %s value %s could not be converted to int64: %s", fieldName, v.Value, err)
					return nil
				}
				return intval
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
