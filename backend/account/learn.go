package account

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
)

var endNumbers = regexp.MustCompile("[0-9]+ *$")
var endStar = regexp.MustCompile("\\*.*$")
var confirmation = regexp.MustCompile(" *[A-Z][0-9][A-Z][0-9][A-Z][0-9] *$")

func (ctx *Context) UpdateLearning(r *Record) (Stats, error) {
	var stats Stats
	if ctx.tx == nil {
		return stats, errors.New("Attempt to update learning without beginning an update")
	}
	if ctx.learnStmt == nil {
		stmt, err := ctx.tx.Prepare(`
insert into learned_cat values(:id, :descr, :amount, :category, :subcategory)
on conflict (pattern, amount)
do update set sourceid=:id, category=:category, subcategory=:subcategory
where :id >= sourceid
                `)
		if err != nil {
			return stats, err
		}
		ctx.learnStmt = stmt
	}

	result, err := ctx.learnStmt.Exec(r.Id, r.Descr, r.Amount, r.Category, r.Subcategory)
	if err != nil {
		return stats, err
	}
	stats.Inserts += getRowsAffected(result)

	noConfirmation := confirmation.ReplaceAllString(r.Descr, "%")
	if noConfirmation != r.Descr {
		result, err = ctx.learnStmt.Exec(r.Id, noConfirmation, r.Amount, r.Category, r.Subcategory)
		if err != nil {
			return stats, err
		}
		stats.Inserts += getRowsAffected(result)
	}

	noNumbers := endNumbers.ReplaceAllString(r.Descr, "%")
	if noNumbers != r.Descr {
		result, err = ctx.learnStmt.Exec(r.Id, noNumbers, r.Amount, r.Category, r.Subcategory)
		if err != nil {
			return stats, err
		}
		stats.Inserts += getRowsAffected(result)
	}

	noStar := endStar.ReplaceAllString(noNumbers, "%")
	if noStar != noNumbers {
		result, err = ctx.learnStmt.Exec(r.Id, noStar, r.Amount, r.Category, r.Subcategory)
		if err != nil {
			return stats, err
		}
		stats.Inserts += getRowsAffected(result)
	}
	return stats, nil
}

func (ctx *Context) ResetLearning() error {
	_, err := ctx.db.Exec(`drop table if exists learned_cat;`)
	if err != nil {
		return err
	}

	_, err = ctx.db.Exec(`
create table learned_cat (
        sourceid integer,
        pattern text,
        amount real,
        category text,
        subcategory text,
        unique(pattern, amount)
);
        `)
	return err
}

func extractCat(row *sql.Row, category *string, subcategory *string) bool {
	var maybeCat sql.NullString
	var maybeSubcat sql.NullString
	err := row.Scan(&maybeCat, &maybeSubcat)
	if err != nil {
		return false
	}
	*category = maybeString(maybeCat)
	*subcategory = maybeString(maybeSubcat)
	return true
}

func (ctx *Context) InferCategory(r *Record) error {
	err, _ := ctx.InferCategoryWithExplanation(r)
	return err
}

func (ctx *Context) InferCategoryWithExplanation(r *Record) (err error, explanation []string) {
	if ctx.findAmtStmt == nil {
		findAmt, err := ctx.db.Prepare(`
select category, subcategory from learned_cat where ? like pattern and amount==?
                `)
		if err != nil {
			return err, explanation
		}
		ctx.findAmtStmt = findAmt
	}

	if ctx.findApproxStmt == nil {
		findApprox, err := ctx.db.Prepare(`
select category, subcategory from learned_cat
where ? like pattern
order by length(pattern) desc, sourceid desc
limit 1
                `)
		if err != nil {
			return err, explanation
		}
		ctx.findApproxStmt = findApprox
	}

	explanation = append(explanation, fmt.Sprintf("looking for descr %s amount %s\n", r.Descr, r.Amount))
	if extractCat(ctx.findAmtStmt.QueryRow(r.Descr, r.Amount), &r.Category, &r.Subcategory) {
		explanation = append(explanation, fmt.Sprintf("descr+amount found %s/%s", r.Category, r.Subcategory))
		return nil, explanation
	}
	explanation = append(explanation, fmt.Sprintf("no descr+amount found; looking for descr %s\n", r.Descr))
	if extractCat(ctx.findApproxStmt.QueryRow(r.Descr), &r.Category, &r.Subcategory) {
		explanation = append(explanation, fmt.Sprintf("descr found %s/%s", r.Category, r.Subcategory))
		return nil, explanation
	}
	explanation = append(explanation, "no descr found; no cat/subcat")
	return nil, explanation
}
