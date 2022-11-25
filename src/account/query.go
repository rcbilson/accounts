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

type QuerySpec struct {
	DateFrom    *Date
	DateUntil   *Date
	DescrLike   *string
	Category    *string
	Subcategory *string
	State       *string
	Limit       *int
	Offset      *int
}

const baseQuery = "SELECT rowid, date, descr, amount, category, subcategory FROM xact"
const orderBy = "ORDER BY date DESC"

func buildQuery(spec QuerySpec) (string, []interface{}) {
	expr := make([]string, 0)
	params := make([]interface{}, 0)
	if spec.DateFrom != nil {
		expr = append(expr, "date >= ?")
		params = append(params, spec.DateFrom.String())
	}
	if spec.DateUntil != nil {
		expr = append(expr, "date < ?")
		params = append(params, spec.DateUntil.String())
	}
	if spec.DescrLike != nil {
		expr = append(expr, "descr like ?")
		params = append(params, *spec.DescrLike)
	}
	if spec.Category != nil {
		expr = append(expr, "category = ?")
		params = append(params, *spec.Category)
	}
	if spec.Subcategory != nil {
		expr = append(expr, "subcategory = ?")
		params = append(params, *spec.Subcategory)
	}
	if spec.State != nil {
		expr = append(expr, "state = ?")
		params = append(params, *spec.State)
	}
	query := baseQuery
	if len(expr) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, expr[0])
		for _, e := range expr[1:] {
			query = fmt.Sprintf("%s AND %s", query, e)
		}
	}
	query = fmt.Sprintf("%s %s", query, orderBy)
	if spec.Limit != nil {
		query = fmt.Sprintf("%s LIMIT ?", query)
		params = append(params, *spec.Limit)
	}
	if spec.Offset != nil {
		query = fmt.Sprintf("%s OFFSET ?", query)
		params = append(params, *spec.Offset)
	}
	return query, params
}

func (ctx *Context) Query(spec QuerySpec) (<-chan *Record, error) {
	query, params := buildQuery(spec)
	rows, err := ctx.db.Query(query, params...)
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
