package learning

import (
	"database/sql"
	"regexp"
)

type LearningContext struct {
	stmt *sql.Stmt
}

func BeginUpdate(tx *sql.Tx) (*LearningContext, error) {
	stmt, err := tx.Prepare(`
insert into learned_cat values(:id, :descr, :amount, :category, :subcategory)
on conflict (pattern, amount)
do update set sourceid=:id, category=:category, subcategory=:subcategory
where :id > sourceid
        `)
	if err != nil {
		return nil, err
	}
	return &LearningContext{stmt}, nil
}

func (learn *LearningContext) EndUpdate() error {
	return nil
}

func getRowsAffected(result sql.Result) int64 {
	rows, _ := result.RowsAffected()
	return rows
}

var endNumbers = regexp.MustCompile("[0-9]+$")
var endStar = regexp.MustCompile("\\*.*$")

func (learn *LearningContext) DoUpdate(id string, descr string, amount string, category string, subcategory string) (int64, error) {
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
