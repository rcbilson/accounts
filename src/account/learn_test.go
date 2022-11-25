package account

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gotest.tools/assert"
)

func setupLearningTest(t *testing.T) *Context {
	os.Setenv("ACCOUNTS_DBFILE", ":memory:")

	acct, err := Open()
	assert.NilError(t, err)
	t.Cleanup(acct.Close)

	err = acct.ResetLearning()
	assert.NilError(t, err)

	return acct
}

func dumpLearning(t *testing.T, acct *Context) {
	rows, err := acct.db.Query("select sourceid, pattern, amount, category, subcategory from learned_cat")
	assert.NilError(t, err)
	defer rows.Close()
	for rows.Next() {
		var sourceid, pattern, amount, category, subcategory string
		err := rows.Scan(&sourceid, &pattern, &amount, &category, &subcategory)
		assert.NilError(t, err)
		fmt.Println(sourceid, pattern, amount, category, subcategory)
	}
}

func TestLearningErrors(t *testing.T) {
	acct := setupLearningTest(t)

	testRecord := Record{"", Date(time.Now()), "Pen Island", "-75.45", "Frivolities", "Tchotchkes"}
	s, err := acct.UpdateLearning(&testRecord)
	assert.Equal(t, s, Stats{0, 0, 0})
	assert.ErrorContains(t, err, "without beginning an update")
}

func TestLearning(t *testing.T) {
	acct := setupLearningTest(t)

	err := acct.BeginUpdate()
	assert.NilError(t, err)
	defer acct.AbortUpdate()

	testRecords := []Record{
		{"1", Date(time.Now()), "SQ * PEN ISLAND # 7545", "-75.45", "Frivolities", "One"},
		{"2", Date(time.Now()), "SQ * PEN ISLAND # 4321", "-43.21", "Frivolities", "Two"},
		{"3", Date(time.Now()), "SQ * QWIK-E-MART # 7545", "-75.45", "Frivolities", "Three"}}
	for _, r := range testRecords {
		s, err := acct.UpdateLearning(&r)
		assert.NilError(t, err)
		assert.Assert(t, s.Inserts > 0)
	}

	err = acct.CompleteUpdate()
	assert.NilError(t, err)

	// exact match
	testRecord := Record{"17", Date(time.Now()), "SQ * PEN ISLAND # 7545", "-75.45", "", ""}
	err = acct.InferCategory(&testRecord)
	assert.NilError(t, err)
	assert.Equal(t, testRecord.Category, "Frivolities")
	assert.Equal(t, testRecord.Subcategory, "One")

	// partial match with exact amount
	testRecord = Record{"17", Date(time.Now()), "SQ * PEN ISLAND # 9989", "-75.45", "", ""}
	err = acct.InferCategory(&testRecord)
	assert.NilError(t, err)
	assert.Equal(t, testRecord.Subcategory, "One")

	// partial match with different amount, should pick most recent option
	testRecord = Record{"17", Date(time.Now()), "SQ * PEN ISLAND # 9989", "-99.89", "", ""}
	err = acct.InferCategory(&testRecord)
	assert.NilError(t, err)
	assert.Equal(t, testRecord.Subcategory, "Two")

	// partial match with exact amount, should pick most recent option
	testRecord = Record{"17", Date(time.Now()), "SQ * AN PIELAND", "-43.21", "", ""}
	err = acct.InferCategory(&testRecord)
	assert.NilError(t, err)
	assert.Equal(t, testRecord.Subcategory, "Two")

	// partial match with different amount, should pick most recent option
	testRecord = Record{"17", Date(time.Now()), "SQ * AN PIELAND", "-98.99", "", ""}
	err = acct.InferCategory(&testRecord)
	assert.NilError(t, err)
	assert.Equal(t, testRecord.Subcategory, "Three")

}
