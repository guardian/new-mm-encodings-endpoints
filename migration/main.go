package main

import (
	"database/sql"
	"flag"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func processTable(db *sql.DB, tableName string) error {
	recordsCh, errCh := AsyncDbReader(db, tableName)

	for {
		select {
		case rec := <-recordsCh:
			if rec == nil {
				log.Print("All done")
				return nil
			}

			log.Printf("DEBUG processTable got %v", rec)
			break
		case err := <-errCh:
			log.Printf("DEBUG processTable got error %s", err)
			return err
		}
	}
}

func main() {
	dsn := flag.String("dsn", "", "MySQL DSN in the form [username[:password]@][protocol[(address)]]/dbname. See https://github.com/go-sql-driver/mysql for details.")
	sourceTable := flag.String("source", "idmapping", "table to read from the SQL database")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("You need to specify -dsn on the commandline. Use --help for more options.")
	}

	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("Could not connect to database at %s: %s", *dsn, err)
	}
	// Recommended settings as per https://github.com/go-sql-driver/mysql
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	processTable(db, *sourceTable)
}
