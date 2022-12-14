package account

import (
	"database/sql"
	"log"
)

type Aggregate struct {
	Category string
	Amount   string
}

func (ctx *Context) aggregate(spec QuerySpec, baseQuery string, orderAndGroup string) (<-chan *Aggregate, error) {
	query, params := buildQuery(spec, baseQuery, orderAndGroup)
	log.Println(query)
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
			err := rows.Scan(&maybeCat, &r.Amount)
			if err != nil {
				log.Println(err)
				return
			}
			r.Category = maybeString(maybeCat)
			log.Println(r)
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
		"GROUP BY category ORDER BY total DESC")
}

func (ctx *Context) AggregateSubcategories(spec QuerySpec) (<-chan *Aggregate, error) {
	return ctx.aggregate(spec,
		"SELECT subcategory, -sum(amount) as total FROM xact",
		"GROUP BY subcategory ORDER BY total DESC")
}
