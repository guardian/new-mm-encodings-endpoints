package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var UrlMatcher = regexp.MustCompile(`^(https?)://[^/]+/(.*)$`)

func makeTargetUrl(endpointBase *string, evt *EndpointEvent) (string, error) {
	matches := UrlMatcher.FindAllStringSubmatch(evt.AccessUrl, -1)
	if matches == nil {
		return "", errors.New(fmt.Sprintf("original URL %s could not be parsed", evt.AccessUrl))
	}
	newUrl := fmt.Sprintf("https://%s/%s", *endpointBase, matches[0][2])
	return newUrl, nil
}

var IrrelevantHeaders = []string{"X-Powered-By"}

func isHeaderIrrelevant(key string) bool {
	for _, hdr := range IrrelevantHeaders {
		if hdr == key {
			return true
		}
	}
	return false
}

func Test(httpClient *http.Client, endpointBase *string, evt *EndpointEvent) (*TestOutput, bool, error) {
	targetUrl, err := makeTargetUrl(endpointBase, evt)
	if err != nil {
		return nil, false, err
	}
	rq, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, false, err
	}

	response, err := httpClient.Do(rq)
	if err != nil {
		return nil, false, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, false, err
	}

	errorList := ""
	success := true
	if response.StatusCode != int(evt.ExpectedResponse) {
		prob := fmt.Sprintf("expected response %d got %d", evt.ExpectedResponse, response.StatusCode)
		errorList += prob + "\n"
		log.Printf("INFO Request %s from %s %s", targetUrl, evt.FormattedTimestamp(), prob)
		success = false
	}
	if string(content) != evt.ExpectedOutputMessage {
		prob := fmt.Sprintf("expected body %s got %s", evt.ExpectedOutputMessage, string(content))
		errorList += prob + "\n"
		log.Printf("INFO Request %s from %s %s", targetUrl, evt.FormattedTimestamp(), prob)
		success = false
	}
	for k, v := range response.Header {
		if headerVal, haveHeader := evt.ExpectedOutputHeaders[k]; haveHeader {
			if headerVal != v[0] {
				prob := fmt.Sprintf("header %s got value %s expected %s", k, v[0], headerVal)
				errorList += prob + "\n"
				log.Printf("INFO Request %s from %s %s", targetUrl, evt.FormattedTimestamp(), prob)
				success = false
			}
		}
	}
	for k, v := range evt.ExpectedOutputHeaders {
		if k == "Content-type" {
			k = "Content-Type"
		}

		if headerVal, haveHeader := response.Header[k]; haveHeader {
			if v == "text/plain" { //fix for an irritation that some events were logged with "text/plain" and others with "text/plain;charset=UTF-8"
				v = "text/plain;charset=UTF-8"
			}
			if headerVal[0] != v {
				prob := fmt.Sprintf("header %s got value %s expected %s", k, headerVal, v)
				errorList += prob + "\n"
				log.Printf("INFO Request %s from %s %s", targetUrl, evt.FormattedTimestamp(), prob)
				success = false
			}
		} else {
			if !isHeaderIrrelevant(k) {
				prob := fmt.Sprintf("response was missing header %s", k)
				errorList += prob + "\n"
				log.Printf("INFO Request %s from %s %s", targetUrl, evt.FormattedTimestamp(), prob)
				success = false
			}
		}
	}

	reformattedHeaders := make(map[string]string, len(response.Header))
	for k, values := range response.Header {
		reformattedHeaders[k] = strings.Join(values, ";")
	}

	ts := time.Now()
	responseEvt := &EndpointEvent{
		Uid:                   uuid.UUID{},
		Timestamp:             &ts,
		AccessUrl:             targetUrl,
		ExpectedOutputMessage: string(content),
		ExpectedOutputHeaders: reformattedHeaders,
		ExpectedResponse:      int16(response.StatusCode),
	}
	out := &TestOutput{
		Request:   evt,
		Result:    responseEvt,
		ErrorList: errorList,
	}
	return out, success, nil
}

func testProcessingThread(inputCh chan *EndpointEvent, outputCh chan TestOutput, endpointBase *string, wg *sync.WaitGroup) {
	successCount := 0
	totalCount := 0

	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	for {
		evt, haveMore := <-inputCh
		if !haveMore {
			log.Printf("AsyncTestEndpoint: reached end of data")
			wg.Done()
			return
		}

		responseEvent, result, err := Test(httpClient, endpointBase, evt)
		if err != nil {
			log.Printf("ERROR Could not perform test for %s at %s: %s", evt.AccessUrl, evt.FormattedTimestamp(), err)
		}
		totalCount++
		if result {
			successCount++
		} else {
			outputCh <- *responseEvent
		}
		log.Printf("INFO Running total %d / %d tests successful", successCount, totalCount)
	}
}

func AsyncTestEndpoint(inputCh chan *EndpointEvent, endpointBase *string, parallel int) (chan TestOutput, *sync.WaitGroup) {
	outputCh := make(chan TestOutput, 100)

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(parallel)

	for i := 0; i < parallel; i++ {
		go testProcessingThread(inputCh, outputCh, endpointBase, waitGroup)
	}

	outputWaitGroup := &sync.WaitGroup{}
	outputWaitGroup.Add(1)
	go func() {
		waitGroup.Wait()
		close(outputCh)
		outputWaitGroup.Done()
	}()

	return outputCh, outputWaitGroup
}
