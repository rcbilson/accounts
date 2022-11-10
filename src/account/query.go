package account

import (
	"database/sql"
	"fmt"
	"log"
)

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
