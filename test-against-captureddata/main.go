package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

func main() {
	tableName := flag.String("table", "", "name of the table to read events from")
	pageSize := flag.Int("s", 50, "page size for event retrieval")
	endpointBase := flag.String("target", "", "server name to test")
	parallel := flag.Int("parallel", 10, "number of requests to run in parallel")
	outputFilename := flag.String("out", "endpoint-test-results.csv", "name of a CSV file to output")
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Could not set up AWS SDK: %s", err)
	}
	ddbClient := dynamodb.NewFromConfig(cfg)

	eventCh, errCh := AsyncRecordReader(ddbClient, tableName, int32(*pageSize))
	//time.Sleep(1*time.Second)
	resultsCh, waitGroup := AsyncTestEndpoint(eventCh, endpointBase, *parallel)

	writeErrCh := AsyncWriter(resultsCh, *outputFilename)
	waitGroup.Add(1)

	go func() {
		for {
			select {
			case err := <-writeErrCh:
				if err == nil {
					waitGroup.Done()
				} else {
					log.Printf("ERROR could not write records: %s", err)
				}
			case err := <-errCh:
				if err != nil {
					log.Printf("ERROR Could not retrieve events: %s", err)
				}
			}
		}
	}()

	log.Print("Waiting for threads to complete...")
	waitGroup.Wait()
	log.Print("Done.")
}
