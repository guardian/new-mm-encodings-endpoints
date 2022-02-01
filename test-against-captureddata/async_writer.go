package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type TestOutput struct {
	Request   *EndpointEvent
	Result    *EndpointEvent
	ErrorList string //newline-separated list of problems detected
}

/*
toCSV returns a list of strings suitable for writing out to CSV
*/
func (t *TestOutput) toCSV() *[]string {
	requestHeaders, _ := json.Marshal(&t.Request.ExpectedOutputHeaders)
	resultHeaders, _ := json.Marshal(&t.Result.ExpectedOutputHeaders)

	log.Printf("DEBUG AsyncWriter got record for %s", t.Request.AccessUrl)
	return &[]string{
		t.Result.AccessUrl,
		t.ErrorList,
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

		writer.Write([]string{"Request URL", "Test result", "Expected status", "Actual status", "Expected output", "Actual output", "Expected headers", "Actual headers"})
		n := 0
		for {
			n++
			evt, moreEvents := <-inputCh
			if !moreEvents {
				log.Printf("INFO AsyncWriter got to end of input, shutting down")
				errCh <- nil
				return
			}

			err := writer.Write(*evt.toCSV())
			writer.Flush()

			if err != nil {
				log.Printf("ERROR Could not write CSV record %d: %s", n, err)
				errCh <- err
				return
			}
		}
	}()
	return errCh
}
