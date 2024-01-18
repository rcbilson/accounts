package learning

import (
	"database/sql"
	"regexp"
)

type Context struct {
	stmt *sql.Stmt
}

func BeginUpdate(tx *sql.Tx) (*Context, error) {
	stmt, err := tx.Prepare(`
insert into learned_cat values(:id, :descr, :amount, :category, :subcategory)
on conflict (pattern, amount)
do update set sourceid=:id, category=:category, subcategory=:subcategory
where :id > sourceid
        `)
	if err != nil {
		return nil, err
	}
	return &Context{stmt}, nil
}

func (learn *Context) EndUpdate() error {
	return nil
}

func getRowsAffected(result sql.Result) int64 {
	rows, _ := result.RowsAffected()
	return rows
}

var endNumbers = regexp.MustCompile("[0-9]+$")
var endStar = regexp.MustCompile("\\*.*$")

func (learn *Context) DoUpdate(id string, descr string, amount string, category string, subcategory string) (int64, error) {
	insertCount := int64(0)

	result, err := learn.stmt.Exec(id, descr, amount, category, subcategory)
	if err != nil {
		return insertCount, err
	}
	insertCount += getRowsAffected(result)

	noNumbers := endNumbers.ReplaceAllString(descr, "%")
	if noNumbers != descr {
		result, err = learn.stmt.Exec(id, noNumbers, amount, category, subcategory)
		if err != nil {
			return insertCount, err
		}
		insertCount += getRowsAffected(result)
	}

	noStar := endStar.ReplaceAllString(noNumbers, "%")
	if noStar != noNumbers {
		result, err = learn.stmt.Exec(id, noStar, amount, category, subcategory)
		if err != nil {
			return insertCount, err
		}
		insertCount += getRowsAffected(result)
	}
	return insertCount, nil
}
