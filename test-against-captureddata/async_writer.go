package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type TestOutput struct {
	Request *EndpointEvent
	Result  *EndpointEvent
}

/*
toCSV returns a list of strings suitable for writing out to CSV
*/
func (t *TestOutput) toCSV() *[]string {
	requestHeaders, _ := json.Marshal(&t.Request.ExpectedOutputHeaders)
	resultHeaders, _ := json.Marshal(&t.Result.ExpectedOutputHeaders)

	return &[]string{
		t.Request.AccessUrl,
		fmt.Sprintf("%d", t.Request.ExpectedResponse),
		fmt.Sprintf("%d", t.Result.ExpectedResponse),
		t.Request.ExpectedOutputMessage,
		t.Result.ExpectedOutputMessage,
		string(requestHeaders),
		string(resultHeaders),
	}
}

/*
AsyncWriter marshals incoming events into CSV rows and writes them to the given output file
*/
func AsyncWriter(inputCh chan TestOutput, filename string) chan error {
	errCh := make(chan error, 1)

	go func() {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
		if err != nil {
			errCh <- err
			return
		}
		defer f.Close()

		writer := csv.NewWriter(f)
		defer writer.Flush()

		n := 0
		for {
			n++
			evt, moreEvents := <-inputCh

			err := writer.Write(*evt.toCSV())
			if err != nil {
				log.Printf("ERROR Could not write CSV record %d: %s", n, err)
				errCh <- err
				return
			}
			if !moreEvents {
				log.Printf("INFO AsyncWriter got to end of input, shutting down")
				errCh <- nil
				return
			}
		}
	}()
	return errCh
}
