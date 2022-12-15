package account

import (
	"database/sql"
	"fmt"
	"log"
)

type Aggregate struct {
	Category string `json:"category"`
	Amount   string `json:"amount"`
}

func (ctx *Context) aggregate(spec QuerySpec, baseQuery string, where string, orderAndGroup string) (<-chan *Aggregate, error) {
	query, params := buildQueryWhere(spec, baseQuery, where, orderAndGroup)
	rows, err := ctx.db.Query(query, params...)
	if err != nil {
		return nil, err
	}

	ch := make(chan *Aggregate)

	go func() {
		defer rows.Close()
		defer close(ch)

		for rows.Next() {
			var r Aggregate
			var maybeCat sql.NullString
			var amount float64
			err := rows.Scan(&maybeCat, &amount)
			if err != nil {
				log.Println(err)
				return
			}
			r.Amount = fmt.Sprintf("%.2f", amount)
			r.Category = maybeString(maybeCat)
			ch <- &r
		}
		err = rows.Err()
		if err != nil {
			log.Println(err)
		}
	}()

	return ch, nil
}

func (ctx *Context) AggregateCategories(spec QuerySpec) (<-chan *Aggregate, error) {
	return ctx.aggregate(spec,
		"SELECT category, -sum(amount) as total FROM xact",
		"category != '' AND category != 'Income'",
		"GROUP BY category ORDER BY total DESC")
}

func (ctx *Context) AggregateSubcategories(spec QuerySpec) (<-chan *Aggregate, error) {
	return ctx.aggregate(spec,
		"SELECT subcategory, -sum(amount) as total FROM xact",
		"",
		"GROUP BY subcategory ORDER BY total DESC")
}
