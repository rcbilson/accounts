package account

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
)

type Context struct {
	db             *sql.DB
	tx             *sql.Tx
	deleteStmt     *sql.Stmt
	updateStmt     *sql.Stmt
	learnStmt      *sql.Stmt
	insertStmt     *sql.Stmt
	findAmtStmt    *sql.Stmt
	findApproxStmt *sql.Stmt
}

type Stats struct {
	Inserts int64
	Deletes int64
	Updates int64
}

func (s *Stats) Add(other Stats) {
	s.Inserts += other.Inserts
	s.Deletes += other.Deletes
	s.Updates += other.Updates
}

func (s Stats) String() string {
	return fmt.Sprintf("%v inserts %v deletes %v updates", s.Inserts, s.Deletes, s.Updates)
}

type specification struct {
	DbFile string
}

func Open() (*Context, error) {
	var ctx Context

	var s specification
	err := envconfig.Process("accounts", &s)
	if err != nil {
		return nil, err
	}

	ctx.db, err = sql.Open("sqlite3", s.DbFile)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}

func (ctx *Context) Close() {
	ctx.db.Close()
}

func (ctx *Context) BeginUpdate() error {
	if ctx.tx != nil {
		return errors.New("Attempt to begin an update while an update is already in progress.")
	}
	var err error
	ctx.tx, err = ctx.db.Begin()
	return err
}

func (ctx *Context) AbortUpdate() {
	if ctx.tx != nil {
		ctx.tx.Rollback()
		ctx.tx = nil
	}
}

func (ctx *Context) CompleteUpdate() error {
	if ctx.tx == nil {
		return errors.New("Attempt to complete an update without an update in progress.")
	}
	err := ctx.tx.Commit()
	ctx.tx = nil
	return err
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
	query := "SELECT rowid, date, descr, amount, category, subcategory FROM xact"
	if where != "" {
		query = fmt.Sprintf("%s WHERE %s", query, where)
	}
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

func getRowsAffected(result sql.Result) int64 {
	rows, _ := result.RowsAffected()
	return rows
}

func (ctx *Context) Delete(id string) (Stats, error) {
	var stats Stats
	if ctx.tx == nil {
		return stats, errors.New("Attempt to delete without beginning an update")
	}
	if ctx.deleteStmt == nil {
		stmt, err := ctx.tx.Prepare("delete from xact where rowid=?")
		if err != nil {
			return stats, err
		}
		ctx.deleteStmt = stmt
	}
	result, err := ctx.deleteStmt.Exec(id)
	if err != nil {
		return stats, err
	}
	stats.Deletes = getRowsAffected(result)
	return stats, nil
}

func (ctx *Context) Update(r *Record) (Stats, error) {
	var stats Stats
	if ctx.tx == nil {
		return stats, errors.New("Attempt to update without beginning an update")
	}
	if r.Id == "" {
		return stats, errors.New("Attempt to update without rowid")
	}
	if ctx.updateStmt == nil {
		stmt, err := ctx.tx.Prepare(
			"update xact set date=?, descr=?, amount=?, category=?, subcategory=?, state=null where rowid=?")
		if err != nil {
			return stats, err
		}
		ctx.updateStmt = stmt
	}
	result, err := ctx.updateStmt.Exec(
		r.Date, r.Descr, r.Amount, r.Category, r.Subcategory, r.Id)
	if err != nil {
		return stats, err
	}
	stats.Updates = getRowsAffected(result)
	s, err := ctx.UpdateLearning(r)
	stats.Add(s)
	return stats, err
}

func (ctx *Context) Insert(r *Record) (Stats, error) {
	var stats Stats
	if ctx.tx == nil {
		return stats, errors.New("Attempt to insert without beginning an update")
	}
	if r.Id != "" {
		return stats, errors.New("Attempt to insert with rowid")
	}
	if ctx.insertStmt == nil {
		insert, err := ctx.tx.Prepare(`
insert into xact (date, descr, amount, category, subcategory, state) values(?, ?, ?, ?, ?, "new")
                `)
		if err != nil {
			return stats, err
		}
		ctx.insertStmt = insert
	}
	result, err := ctx.insertStmt.Exec(
		r.Date, r.Descr, r.Amount, r.Category, r.Subcategory)
	if err != nil {
		return stats, err
	}
	stats.Updates = getRowsAffected(result)
	return stats, err
}
