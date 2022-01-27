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
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("Could not set up AWS SDK: %s", err)
	}
	ddbClient := dynamodb.NewFromConfig(cfg)

	eventCh, errCh := AsyncRecordReader(ddbClient, tableName, int32(*pageSize))

	for {
		select {
		case event, moreEvents := <-eventCh:
			log.Printf("Got event: %v", event)
			if !moreEvents {
				log.Printf("All done!")
				return
			}
		case err := <-errCh:
			log.Printf("ERROR Could not retrieve events: %s", err)
		}
	}
}
