package main

import (
	"database/sql"
	"log"
	"regexp"

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

var endNumbers = regexp.MustCompile("[0-9]+$")
var endStar = regexp.MustCompile("\\*.*$")

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

	stmt, err := tx.Prepare(`
insert into learned_cat values(:id, :descr, :amount, :category, :subcategory)
on conflict (pattern, amount)
do update set sourceid=:id, category=:category, subcategory=:subcategory
where :id > sourceid
        `)
	check(err)

	rowCount := 0
	insertCount := int64(0)
	for _, r := range records {
		rowCount++

		if r.id == "" {
			panic("aarg")
		}

		result, err := stmt.Exec(r.id, r.descr, r.amount, r.maybeCat, r.maybeSubcat)
		checkRow(err, rowCount)
		insertCount += getRowsAffected(result)

		noNumbers := endNumbers.ReplaceAllString(r.descr, "%")
		if noNumbers != r.descr {
			result, err = stmt.Exec(r.id, noNumbers, r.amount, r.maybeCat, r.maybeSubcat)
			checkRow(err, rowCount)
			insertCount += getRowsAffected(result)
		}

		noStar := endStar.ReplaceAllString(noNumbers, "%")
		if noStar != noNumbers {
			result, err = stmt.Exec(r.id, noStar, r.amount, r.maybeCat, r.maybeSubcat)
			checkRow(err, rowCount)
			insertCount += getRowsAffected(result)
		}
	}
	err = tx.Commit()
	check(err)

	log.Println(rowCount, "rows processed", insertCount, "inserted")
}
