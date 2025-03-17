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

type MonthAggregate struct {
	Aggregate
	Month string `json:"month"`
}

type SummaryChart struct {
	Amounts []MonthAggregate `json:"amounts"`
}

func buildMonthAggregate(rows *sql.Rows) (<-chan *MonthAggregate, error) {
	ch := make(chan *MonthAggregate)

	go func() {
		defer rows.Close()
		defer close(ch)

		for rows.Next() {
			var r MonthAggregate
			var maybeCat sql.NullString
			var amount float64
			err := rows.Scan(&r.Month, &maybeCat, &amount)
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

// as a percentage of income
//
//with incomebymonth as
//	(select strftime("%Y-%m", date) as mth, sum(amount) as inc from xact where category="Income" group by mth)
//select strftime("%Y-%m", date) as month, kc.key, sum(-amount)*100/ibm.inc from xact
//join keycats as kc on xact.category=kc.category
//join incomebymonth as ibm on month=ibm.mth

// monthly total
//
//select strftime("%Y-%m", date) as month, kc.key, sum(-amount) from xact
//join keycats as kc on xact.category=kc.category

func (ctx *Context) SummaryChart(spec QuerySpec) (*SummaryChart, error) {
	var result SummaryChart

	dateWhere := ""
	if spec.DateFrom != nil {
		dateWhere = fmt.Sprintf("where date >= '%s'", spec.DateFrom.String())
	}
	if spec.DateUntil != nil {
		if dateWhere != "" {
			dateWhere = fmt.Sprintf("%s AND date <= '%s'", dateWhere, spec.DateUntil.String())
		} else {
			dateWhere = fmt.Sprintf("where date <= '%s'", spec.DateUntil.String())
		}
	}

	query := fmt.Sprintf(
		`
with months as (
	select distinct strftime("%%Y-%%m", date) as month from xact %[1]s
),
keys as (
	select distinct key from keycats
),
month_keys as (
	select month, key from months cross join keys
),
month_totals as (
	select
		strftime("%%Y-%%m", date) as month,
		kc.key,
		sum(-amount) as total
	from xact
	join keycats as kc on xact.category=kc.category
	%[1]s
	group by month, kc.key
),
month_filled as (
	select
		mk.month,
		mk.key,
		coalesce(mt.total, 0) as total
	from month_keys mk
	left join month_totals mt on mk.month = mt.month and mk.key = mt.key
),
month_cumulative as (
    select
		month,
		key,
		sum(total) over (partition by key order by month rows between unbounded preceding and current row) as total
	from month_filled mf
),
income_totals as (
        select strftime("%%Y-%%m", date) as month, sum(amount) as total from xact where category="Income" group by month
),
income_filled as (
        select
                mk.month,
                coalesce(mt.total, 0) as total
        from months mk
        left join income_totals mt on mk.month = mt.month
),
income_cumulative as (
        select
                month,
                sum(total) over (order by month rows between unbounded preceding and current row) as total
        from income_filled mf
)
select
        mf.month,
        mf.key,
        mf.total*100/ic.total as total
from month_cumulative mf
join income_cumulative ic on mf.month=ic.month
order by mf.month asc
		`,
		dateWhere)
	rows, err := ctx.db.Query(query)
	if err != nil {
		return nil, err
	}

	ch, err := buildMonthAggregate(rows)
	if err != nil {
		return nil, err
	}

	result.Amounts = make([]MonthAggregate, 0)
	for agg := range ch {
		result.Amounts = append(result.Amounts, *agg)
	}

	return &result, nil
}
