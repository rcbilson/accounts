package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
)

type Specification struct {
	DbFile string
}

func check(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func checkRow(err error, rowCount int) {
	if err != nil {
		log.Fatalf("at line %d: %v", rowCount, err.Error())
	}
}

func getRowsAffected(result sql.Result) int64 {
	rows, err := result.RowsAffected()
	check(err)
	return rows
}

func doDelete(db *sql.DB, id string, rowCount int) int64 {
	result, err := db.Exec("delete from xact where rowid=?", id)
	checkRow(err, rowCount)
	return getRowsAffected(result)
}

func doUpdate(db *sql.DB, record []string, rowCount int) int64 {
	id := record[0]
	date, err := time.Parse("2006-01-02", record[1])
	checkRow(err, rowCount)
	descr := record[2]
	amount, err := strconv.ParseFloat(record[3], 64)
	checkRow(err, rowCount)
	category := record[4]
	subcategory := record[5]
	result, err := db.Exec(
		"update xact set date=?, descr=?, amount=?, category=?, subcategory=?, state=null where rowid=?",
		date, descr, amount, category, subcategory, id)
	checkRow(err, rowCount)
	return getRowsAffected(result)
}

func main() {
	var s Specification
	err := envconfig.Process("accounts", &s)
	check(err)

	db, err := sql.Open("sqlite3", s.DbFile)
	if err != nil {
		log.Fatal(err)
	}

	csvFile := os.Args[1]
	f, err := os.Open(csvFile)
	check(err)

	r := csv.NewReader(f)
	// Don't insist on all records having the same number of fields.
	// This allows us to use a bare rowid to indicate a deletion.
	// It does mean that we have to manually verify that the correct
	// number of fields are present.
	r.FieldsPerRecord = -1

	rowCount := 0
	deleteCount := int64(0)
	updateCount := int64(0)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		rowCount++
		checkRow(err, rowCount)
		if len(record) == 1 {
			deleteCount += doDelete(db, record[0], rowCount)
		} else if len(record) == 6 {
			updateCount += doUpdate(db, record, rowCount)
		} else {
			checkRow(errors.New("unexpected number of fields in row"), rowCount)
		}
	}
	log.Println(rowCount, "rows processed", updateCount, "updates", deleteCount, "deletions")
}
