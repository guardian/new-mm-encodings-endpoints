package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
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

func Test(httpClient *http.Client, endpointBase *string, evt *EndpointEvent) (bool, error) {
	targetUrl, err := makeTargetUrl(endpointBase, evt)
	if err != nil {
		return false, err
	}
	rq, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return false, err
	}

	response, err := httpClient.Do(rq)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	success := true
	if response.StatusCode != int(evt.ExpectedResponse) {
		log.Printf("INFO Request %s from %s expected response %d got %d", evt.AccessUrl, evt.FormattedTimestamp(), evt.ExpectedResponse, response.StatusCode)
		success = false
	}
	if string(content) != evt.ExpectedOutputMessage {
		log.Printf("INFO Request %s from %s expected body %s got %s", evt.AccessUrl, evt.FormattedTimestamp(), evt.ExpectedOutputMessage, string(content))
		success = false
	}
	for k, v := range response.Header {
		if headerVal, haveHeader := evt.ExpectedOutputHeaders[k]; haveHeader {
			if headerVal != v[0] {
				log.Printf("INFO Request %s from %s header %s got value %s expected %s", evt.AccessUrl, evt.FormattedTimestamp(), k, v[0], headerVal)
				success = false
			}
		} else {
			//log.Printf("INFO Request %s from %s response had extra header %s", evt.AccessUrl, evt.FormattedTimestamp(), k)
		}
	}
	for k, v := range evt.ExpectedOutputHeaders {
		if headerVal, haveHeader := response.Header[k]; haveHeader {
			if headerVal[0] != v {
				log.Printf("INFO Request %s from %s header %s got value %s expected %s", evt.AccessUrl, evt.FormattedTimestamp(), k, v[0], headerVal)
				success = false
			}
		} else {
			log.Printf("INFO Request %s from %s response was missing header %s", evt.AccessUrl, evt.FormattedTimestamp(), k)
			success = false
		}
	}
	return success, nil
}

func testProcessingThread(inputCh chan *EndpointEvent, endpointBase *string, wg *sync.WaitGroup) {
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
			errCh <- nil
			return
		}

		result, err := Test(httpClient, endpointBase, evt)
		if err != nil {
			log.Printf("ERROR Could not perform test for %s at %s: %s", evt.AccessUrl, evt.FormattedTimestamp(), err)
		}
		totalCount++
		if result {
			successCount++
		}
		log.Printf("INFO Running total %d / %d tests successful", successCount, totalCount)
	}
}
func AsyncTestEndpoint(inputCh chan *EndpointEvent, endpointBase *string, parallel int) (chan TestOutput, chan error) {
	outputCh := make(chan TestOutput, 100)
	errCh := make(chan error, 1)

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(parallel)

	return outputCh, errCh
}
