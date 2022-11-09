package main

import (
	"database/sql"
	"encoding/csv"
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

func maybeString(maybe sql.NullString) string {
	if maybe.Valid {
		return maybe.String
	} else {
		return ""
	}
}

type importRecord struct {
	date   time.Time
	descr  string
	amount float64
}

func parseTD(record []string) (*importRecord, error) {
	// 01/05/2021,SEND E-TFR CA***t5J ,413.00,,13041.20
	var err error
	var r importRecord
	r.date, err = time.Parse("01/02/2006", record[0])
	if err != nil {
		return nil, err
	}
	r.descr = record[1]
	debit, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		debit = 0
	}
	credit, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		credit = 0
	}
	r.amount = credit - debit
	return &r, nil
}

func parseTangerine(record []string) (*importRecord, error) {
	// 1/1/2021,DEBIT,Recurring Internet Withdrawal to,SPOINT To 3023010786,-250
	var err error
	var r importRecord
	r.date, err = time.Parse("1/2/2006", record[0])
	if err != nil {
		return nil, err
	}
	r.descr = record[2] + "/" + record[3]
	r.amount, err = strconv.ParseFloat(record[4], 64)
	if err != nil {
		r.amount = 0
	}
	return &r, nil
}

func extractCat(row *sql.Row, category *string, subcategory *string) bool {
	var maybeCat sql.NullString
	var maybeSubcat sql.NullString
	err := row.Scan(&maybeCat, &maybeSubcat)
	if err == sql.ErrNoRows {
		return false
	}
	check(err)
	*category = maybeString(maybeCat)
	*subcategory = maybeString(maybeSubcat)
	return true
}

func main() {
	var s Specification
	err := envconfig.Process("accounts", &s)
	check(err)

	db, err := sql.Open("sqlite3", s.DbFile)
	check(err)

	findAmt, err := db.Prepare(`
select category, subcategory from learned_cat where ? like pattern and amount==?
        `)
	check(err)

	findApprox, err := db.Prepare(`
select category, subcategory from learned_cat where ? like pattern order by length(pattern), sourceid desc limit 1
        `)
	check(err)

	tx, err := db.Begin()
	check(err)
	defer tx.Rollback()

	insert, err := tx.Prepare(`
insert into xact (date, descr, amount, category, subcategory, state) values(?, ?, ?, ?, ?, "new")
        `)
	check(err)

	for _, csvFile := range os.Args[1:] {
		f, err := os.Open(csvFile)
		check(err)

		r := csv.NewReader(f)

		parser := parseTD
		rowCount := 0
		insertCount := int64(0)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			rowCount++
			if record[0] == "Date" {
				parser = parseTangerine
				continue
			}
			rec, err := parser(record)
			if err != nil {
				log.Fatalf("at %s row %d: %v", csvFile, rowCount, err)
			}
			var category string
			var subcategory string
			if !extractCat(findAmt.QueryRow(rec.descr, rec.amount), &category, &subcategory) {
				extractCat(findApprox.QueryRow(rec.descr), &category, &subcategory)
			}

			result, err := insert.Exec(rec.date, rec.descr, rec.amount, category, subcategory)
			checkRow(err, rowCount)
			rows, err := result.RowsAffected()
			check(err)
			insertCount += rows
		}

		log.Printf("%s: %v rows processed %v inserted", csvFile, rowCount, insertCount)
	}
	err = tx.Commit()
	check(err)
}
