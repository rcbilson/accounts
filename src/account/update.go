package account

import (
	"database/sql"
	"errors"
	"fmt"
)

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
