package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
	"knilson.org/accounts/learning"
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

func doDelete(tx *sql.Tx, id string, rowCount int) int64 {
	result, err := tx.Exec("delete from xact where rowid=?", id)
	checkRow(err, rowCount)
	return getRowsAffected(result)
}

func doUpdate(tx *sql.Tx, learn *learning.Context, record []string, rowCount int) int64 {
	id := record[0]
	date, err := time.Parse("2006-01-02", record[1])
	checkRow(err, rowCount)
	descr := record[2]
	amount, err := strconv.ParseFloat(record[3], 64)
	checkRow(err, rowCount)
	category := ""
	subcategory := ""
	if len(record) > 4 {
		category = record[4]
	}
	if len(record) > 5 {
		subcategory = record[5]
	}
	result, err := tx.Exec(
		"update xact set date=?, descr=?, amount=?, category=?, subcategory=?, state=null where rowid=?",
		date, descr, amount, category, subcategory, id)
	checkRow(err, rowCount)
	_, err = learn.DoUpdate(id, descr, fmt.Sprintf("%f", amount), category, subcategory)
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
	defer db.Close()

	tx, err := db.Begin()
	check(err)
	defer tx.Rollback()

	learn, err := learning.BeginUpdate(tx)
	check(err)
	defer learn.EndUpdate()

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
			deleteCount += doDelete(tx, record[0], rowCount)
		} else if len(record) >= 4 {
			updateCount += doUpdate(tx, learn, record, rowCount)
		} else {
			checkRow(errors.New("unexpected number of fields in row"), rowCount)
		}
	}
	err = tx.Commit()
	check(err)
	log.Println(rowCount, "rows processed", updateCount, "updates", deleteCount, "deletions")
}
