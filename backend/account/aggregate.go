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

	return buildAggregate(rows)
}

func buildAggregate(rows *sql.Rows) (<-chan *Aggregate, error) {
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
		err := rows.Err()
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

type Summary struct {
	Income  string      `json:"income"`
	Amounts []Aggregate `json:"amounts"`
}

func (ctx *Context) Summary(spec QuerySpec) (*Summary, error) {
	var result Summary

	query, params := buildQueryWhere(spec,
		"select sum(amount) from xact",
		"category = 'Income'",
		"")
	row := ctx.db.QueryRow(query, params...)
	var income float64
	err := row.Scan(&income)
	if err != nil {
		return nil, err
	}
	result.Income = fmt.Sprintf("%.2f", income)

	query, params = buildQueryWhere(spec,
		"SELECT key, -sum(amount) AS total FROM xact JOIN keycats ON keycats.category = xact.category",
		"",
		"GROUP BY key ORDER BY total DESC")
	rows, err := ctx.db.Query(query, params...)
	if err != nil {
		return nil, err
	}

	ch, err := buildAggregate(rows)
	if err != nil {
		return nil, err
	}

	result.Amounts = make([]Aggregate, 0)
	for agg := range ch {
		result.Amounts = append(result.Amounts, *agg)
	}

	return &result, nil
}
