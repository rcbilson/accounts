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

func main() {
	var s Specification
	err := envconfig.Process("accounts", &s)
	check(err)

	db, err := sql.Open("sqlite3", s.DbFile)
	check(err)

	_, err = db.Exec(`
create temporary table imported (
	date date,
	descr text,
	amount real
);
        `)
	check(err)

	for _, csvFile := range os.Args[1:] {
		f, err := os.Open(csvFile)
		check(err)

		r := csv.NewReader(f)

		parser := parseTD
		rowCount := 0
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
			_, err = db.Exec("insert into imported(date, descr, amount) values (?, ?, ?);",
				rec.date, rec.descr, rec.amount)
		}
	}

	_, err = db.Exec(`
create temporary view possibilities as select imported.rowid as id, imported.date, imported.descr, imported.amount, catmap.category, catmap.subcategory, catmap.score from imported left outer join catmap on imported.descr like catmap.descr;
        `)
	check(err)
	result, err := db.Exec(`
insert into xact select date, descr, amount, category, subcategory, "new" from possibilities join (select id as bestid, max(score) as bestscore from possibilities group by id) where bestid = id and (bestscore is null or score = bestscore);
        `)
	check(err)
	rows, err := result.RowsAffected()
	check(err)
	log.Println(rows, "rows inserted")
}
