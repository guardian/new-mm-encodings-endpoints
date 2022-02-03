package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

/*
toCSV returns a list of strings suitable for writing out to CSV
*/
func (t *TestOutput) toJsonBytes() *[]byte {
	log.Printf("DEBUG AsyncJsonWriter got record for %s", t.Request.AccessUrl)
	//return &[]string{
	//	t.Result.AccessUrl,
	//	t.ErrorList,
	//	fmt.Sprintf("%d", t.Request.ExpectedResponse),
	//	fmt.Sprintf("%d", t.Result.ExpectedResponse),
	//	t.Request.ExpectedOutputMessage,
	//	t.Result.ExpectedOutputMessage,
	//	string(requestHeaders),
	//	string(resultHeaders),
	//}
	rawContent := map[string]interface{}{
		"testedUrl":        t.Result.AccessUrl,
		"errorList":        strings.Split(t.ErrorList, "\n"),
		"expectedResponse": t.Request.ExpectedResponse,
		"actualResponse":   t.Result.ExpectedResponse,
	}
	if expectedLoc, haveLoc := t.Request.ExpectedOutputHeaders["Location"]; haveLoc {
		rawContent["expectedLocation"] = expectedLoc
		rawContent["actualLocation"] = t.Result.ExpectedOutputHeaders["Location"]
	}
	content, err := json.Marshal(&rawContent)
	if err != nil {
		log.Fatalf("Could not marshal %v: %s", rawContent, err)
	}
	return &content
}

/*
AsyncWriter marshals incoming events into CSV rows and writes them to the given output file
*/
func AsyncJsonWriter(inputCh chan TestOutput, filename string) chan error {
	errCh := make(chan error, 1)

	go func() {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
		if err != nil {
			errCh <- err
			return
		}

		f.WriteString("[\n")
		defer func() {
			f.WriteString("\n]\n")
			f.Close()
		}()

		n := 0
		for {
			n++
			evt, moreEvents := <-inputCh
			if !moreEvents {
				log.Printf("INFO AsyncWriter got to end of input, shutting down")
				errCh <- nil
				return
			}

			if n > 1 {
				f.WriteString(",\n")
			}
			_, err := f.Write(*evt.toJsonBytes())
			if err != nil {
				log.Printf("ERROR Could not write JSON record %d: %s", n, err)
				errCh <- err
				return
			}
		}
	}()
	return errCh
}
