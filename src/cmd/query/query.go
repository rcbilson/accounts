package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
)

func maybeString(maybe sql.NullString) string {
	if maybe.Valid {
		return maybe.String
	} else {
		return ""
	}
}

func query2csv(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	var id string
	var date time.Time
	var descr string
	var amount string
	var category string
	var subcategory string
	var maybeCat sql.NullString
	var maybeSubcat sql.NullString
	for rows.Next() {
		err := rows.Scan(&id, &date, &descr, &amount, &maybeCat, &maybeSubcat)
		if err != nil {
			log.Fatal(err)
		}
		category = maybeString(maybeCat)
		subcategory = maybeString(maybeSubcat)
		record := []string{id, date.Format("2006-01-02"), descr, amount, category, subcategory}
		writer.Write(record)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

type Specification struct {
	DbFile string
}

func main() {
	var s Specification
	err := envconfig.Process("accounts", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := sql.Open("sqlite3", s.DbFile)
	if err != nil {
		log.Fatal(err)
	}

	where := os.Args[1]
	q := fmt.Sprintf("SELECT rowid, date, descr, amount, category, subcategory FROM xact WHERE %s", where)

	query2csv(db, q)
}
