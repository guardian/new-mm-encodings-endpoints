package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func processTable(db *sql.DB, tableName string, ddbClient *dynamodb.Client, outputTableName *string, addUuid bool, nullableKeyFields string) error {
	recordsCh, errCh := AsyncDbReader(db, tableName)
	writeErrCh := AsyncDynamoWriter(recordsCh, ddbClient, outputTableName, addUuid, nullableKeyFields)

	for {
		select {
		//case rec := <-recordsCh:
		//	if rec == nil {
		//		log.Print("All done")
		//		return nil
		//	}
		//
		//	log.Printf("DEBUG processTable got %v", rec)
		//	break
		case err := <-writeErrCh:
			if err == nil {
				log.Print("All done")
				return nil
			}
			log.Printf("ERROR processTable got error %s", err)
			return err
		case err := <-errCh:
			log.Printf("DEBUG processTable got error %s", err)
			return err
		}
	}
}

func main() {
	dsn := flag.String("dsn", "", "MySQL DSN in the form [username[:password]@][protocol[(address)]]/dbname. See https://github.com/go-sql-driver/mysql for details.")
	sourceTable := flag.String("source", "idmapping", "table to read from the SQL database")
	destTable := flag.String("dest", "", "dynamodb table to write to")
	addUUID := flag.Bool("add-uuid", false, "add a uniquely generated id if this is specified")
	nullableKeyFields := flag.String("nullable-fields", "", "A comma separated list of fields that can be null")
	flag.Parse()

	uuid.EnableRandPool()

	if *dsn == "" {
		log.Fatal("You need to specify -dsn on the commandline. Use --help for more options.")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ddbClient := dynamodb.NewFromConfig(cfg)

	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("Could not connect to database at %s: %s", *dsn, err)
	}
	// Recommended settings as per https://github.com/go-sql-driver/mysql
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = processTable(db, *sourceTable, ddbClient, destTable, *addUUID, *nullableKeyFields)
	if err != nil {
		log.Fatal("Error exit")
	} else {
		os.Exit(0)
	}
}
