package account

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
)

type Context struct {
	db *sql.DB
	tx *sql.Tx
}

type specification struct {
	DbFile string
}

func Open() (*Context, error) {
	var s specification
	err := envconfig.Process("accounts", &s)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", s.DbFile)
	if err != nil {
		return nil, err
	}

	return &Context{db, nil}, nil
}

func (ctx *Context) BeginUpdate() {
}

func (ctx *Context) AbortUpdate() {
}

func (ctx *Context) CompleteUpdate() {
}

type Record struct {
	Id          string
	Date        time.Time
	Descr       string
	Amount      string
	Category    string
	Subcategory string
}

func maybeString(maybe sql.NullString) string {
	if maybe.Valid {
		return maybe.String
	} else {
		return ""
	}
}

func (ctx *Context) Query(where string) (<-chan *Record, error) {
	query := fmt.Sprintf("SELECT rowid, date, descr, amount, category, subcategory FROM xact WHERE %s", where)
	rows, err := ctx.db.Query(query)
	if err != nil {
		return nil, err
	}

	ch := make(chan *Record)

	go func() {
		defer rows.Close()
		defer close(ch)

		for rows.Next() {
			var r Record
			var maybeCat sql.NullString
			var maybeSubcat sql.NullString
			err := rows.Scan(&r.Id, &r.Date, &r.Descr, &r.Amount, &maybeCat, &maybeSubcat)
			if err != nil {
				log.Println(err)
				return
			}
			r.Category = maybeString(maybeCat)
			r.Subcategory = maybeString(maybeSubcat)
			ch <- &r
		}
		err = rows.Err()
		if err != nil {
			log.Println(err)
		}
	}()

	return ch, nil
}
