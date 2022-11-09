package main

import (
	"database/sql"
	"log"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
	"knilson.org/accounts/learning"
)

type Specification struct {
	DbFile string
}

func maybeString(maybe sql.NullString) string {
	if maybe.Valid {
		return maybe.String
	} else {
		return ""
	}
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

func main() {
	var s Specification
	err := envconfig.Process("accounts", &s)
	check(err)

	db, err := sql.Open("sqlite3", "file:"+s.DbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	_, err = db.Exec(`drop table if exists learned_cat;`)
	check(err)

	_, err = db.Exec(`
create table learned_cat (
        sourceid integer,
        pattern text,
        amount real,
        category text,
        subcategory text,
        unique(pattern, amount)
);
        `)
	check(err)

	type record struct {
		id          string
		descr       string
		amount      string
		maybeCat    sql.NullString
		maybeSubcat sql.NullString
	}
	records := make([]record, 0)

	{
		rows, err := db.Query("select rowid, descr, amount, category, subcategory from xact;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var r record
			err := rows.Scan(&r.id, &r.descr, &r.amount, &r.maybeCat, &r.maybeSubcat)
			if err != nil {
				log.Fatal(err)
			}
			records = append(records, r)
		}
	}

	tx, err := db.Begin()
	check(err)
	defer tx.Rollback()

	learn, err := learning.BeginUpdate(tx)
	check(err)

	rowCount := 0
	insertCount := int64(0)
	for _, r := range records {
		rowCount++
		inserted, err := learn.DoUpdate(r.id, r.descr, r.amount, maybeString(r.maybeCat), maybeString(r.maybeSubcat))
		checkRow(err, rowCount)
		insertCount += inserted
	}
	err = tx.Commit()
	check(err)

	log.Println(rowCount, "rows processed", insertCount, "inserted")
}
