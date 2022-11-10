package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"knilson.org/accounts/account"
)

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

func doDelete(acct *account.Context, id string, rowCount int) account.Stats {
	stats, err := acct.Delete(id)
	checkRow(err, rowCount)
	return stats
}

func doUpdate(acct *account.Context, record []string, rowCount int) account.Stats {
	var r account.Record
	var err error
	r.Id = record[0]
	r.Date, err = time.Parse("2006-01-02", record[1])
	checkRow(err, rowCount)
	r.Descr = record[2]
	_, err = strconv.ParseFloat(record[3], 64)
	checkRow(err, rowCount)
	r.Amount = record[3]
	if len(record) > 4 {
		r.Category = record[4]
	}
	if len(record) > 5 {
		r.Subcategory = record[5]
	}
	stats, err := acct.Update(&r)
	checkRow(err, rowCount)
	return stats
}

func main() {
	acct, err := account.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer acct.Close()

	err = acct.BeginUpdate()
	check(err)
	defer acct.AbortUpdate()

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
	var stats account.Stats
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		rowCount++
		checkRow(err, rowCount)
		var s account.Stats
		if len(record) == 1 {
			s = doDelete(acct, record[0], rowCount)
		} else if len(record) >= 4 {
			s = doUpdate(acct, record, rowCount)
		} else {
			checkRow(errors.New("unexpected number of fields in row"), rowCount)
		}
		stats.Add(s)
	}
	err = acct.CompleteUpdate()
	check(err)
	log.Println(rowCount, "rows processed", stats)
}
