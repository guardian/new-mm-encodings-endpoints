package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type GeneralRecord map[string]interface{}

/**
takes an untyped byte array from mysql and converts it into a string (assuming utf-8 bytes)
*/
func derefTypeString(value *interface{}) string {
	if value == nil || *value == nil {
		return ""
	}

	byteValue := (*value).([]byte)
	return string(byteValue)
}

/*
AsyncDbReader scans an entire mysql table and outputs generic (typed) records of map[string]interface{}.
The output stream terminates with a `nil` value if successful, or a single value in the error channel if unsuccessful.

Columns are converted to Go native data types before being output to the map.
*/
func AsyncDbReader(db *sql.DB, tableToScan string) (chan GeneralRecord, chan error) {
	outputCh := make(chan GeneralRecord, 100)
	errCh := make(chan error, 1)

	go func() {
		q := fmt.Sprintf("select * from %s", tableToScan)

		rowsPtr, err := db.Query(q)
		if err != nil {
			log.Printf("ERROR AsyncDbReader could not query %s: %s", tableToScan, err)
			errCh <- err
			return
		}
		defer rowsPtr.Close()

		columns, err := rowsPtr.Columns()
		if err != nil {
			errCh <- err
			return
		}

		colTypes, err := rowsPtr.ColumnTypes()
		if err != nil {
			errCh <- err
			return
		}

		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		for rowsPtr.Next() {
			rec := make(GeneralRecord, len(columns))
			err := rowsPtr.Scan(values...)
			if err != nil {
				log.Printf("ERROR AsyncDbReader could not scan data from %s: %s", tableToScan, err)
				errCh <- err
				return
			}

			for i, col := range columns {
				switch colTypes[i].DatabaseTypeName() {
				case "TINYINT": //8 bit integer
					stringValue := derefTypeString(values[i].(*interface{}))
					if stringValue != "" {
						bigint, err := strconv.ParseInt(stringValue, 10, 8)
						if err != nil {
							log.Printf("WARNING could not convert TINYINT value %s: %s", stringValue, err)
						}
						rec[col] = int8(bigint)
					}
					break
				case "SMALLINT": //16 bit integer
					stringValue := derefTypeString(values[i].(*interface{}))
					if stringValue != "" {
						bigint, err := strconv.ParseInt(stringValue, 10, 16)
						if err != nil {
							log.Printf("WARNING could not convert SMALLINT value %s: %s", stringValue, err)
						}
						rec[col] = int16(bigint)
					}
					break
				case "INT": //32 bit integer
					stringValue := derefTypeString(values[i].(*interface{}))
					if stringValue != "" {
						bigint, err := strconv.ParseInt(stringValue, 10, 32)
						if err != nil {
							log.Printf("WARNING could not convert INT value %s: %s", stringValue, err)
						}
						rec[col] = int32(bigint)
					}
					break
				case "BIGINT": //64 bit integer
					stringValue := derefTypeString(values[i].(*interface{}))
					if stringValue != "" {
						rec[col], err = strconv.ParseInt(stringValue, 10, 64)
						if err != nil {
							log.Printf("WARNING could not convert INT value %s: %s", stringValue, err)
						}
					}
					break
				case "TIMESTAMP":
					//stringValue := *(values[i].(*string))
					stringValue := derefTypeString(values[i].(*interface{}))

					timeValue, err := time.Parse("2006-01-02 15:04:05", stringValue)
					if err != nil {
						log.Printf("WARNING invalid time value %s: %s", stringValue, err)
					} else {
						rec[col] = timeValue
					}
					break
				case "FLOAT":
					stringValue := derefTypeString(values[i].(*interface{}))
					if stringValue != "" {
						floatValue, err := strconv.ParseFloat(stringValue, 64)
						if err != nil {
							log.Printf("WARNING invalid floag value %s: %s", floatValue, err)
						} else {
							rec[col] = floatValue
						}
					}
				//	rec[col] = t
				//case float64:
				//	rec[col] = t
				case "VARCHAR":
					fallthrough
				case "TEXT":
					//rec[col] = *(values[i].(*string))
					rec[col] = derefTypeString(values[i].(*interface{}))
					break
				//case bool:
				//	rec[col] = t
				default:
					log.Printf("WARNING unrecognised type %s", colTypes[i].DatabaseTypeName())
				}
			}

			outputCh <- rec
		}
		outputCh <- nil
	}()

	return outputCh, errCh
}
