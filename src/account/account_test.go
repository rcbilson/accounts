package account

import (
	"os"
	"testing"
	"time"

	"gotest.tools/assert"
)

func setupTest(t *testing.T) *Context {
	os.Setenv("ACCOUNTS_DBFILE", ":memory:")

	acct, err := Open()
	assert.NilError(t, err)
	t.Cleanup(acct.Close)

	_, err = acct.db.Exec(`
CREATE TABLE xact (
        date date,
        descr text,
        amount real,
        category text,
        subcategory text,
        state text);
        `)
	assert.NilError(t, err)

	return acct
}

func materializeQuery(t *testing.T, acct *Context, query QuerySpec) []Record {
	ch, err := acct.Query(query)
	assert.NilError(t, err)

	result := make([]Record, 0)
	for rec := range ch {
		result = append(result, *rec)
	}
	return result
}

func assertRecordsEqualNoId(t *testing.T, left Record, right Record) {
	assert.Equal(t, left.Descr, right.Descr)
	assert.Equal(t, left.Amount, right.Amount)
	assert.Equal(t, left.Category, right.Category)
	assert.Equal(t, left.Subcategory, right.Subcategory)
	assert.Equal(t, left.Date.String(), right.Date.String())
}

func assertRecordsEqual(t *testing.T, left Record, right Record) {
	assert.Equal(t, left.Id, right.Id)
	assertRecordsEqualNoId(t, left, right)
}

func testInsertion(t *testing.T, acct *Context, r *Record) {
	err := acct.BeginUpdate()
	assert.NilError(t, err)
	defer acct.AbortUpdate()

	s, err := acct.Insert(r)
	assert.NilError(t, err)

	assert.Equal(t, s, Stats{1, 0, 0})

	err = acct.CompleteUpdate()
	assert.NilError(t, err)

}

func testUpdate(t *testing.T, acct *Context, r *Record) {
	err := acct.BeginUpdate()
	assert.NilError(t, err)
	defer acct.AbortUpdate()

	s, err := acct.Update(r)
	assert.NilError(t, err)

	assert.Equal(t, s, Stats{0, 0, 1})

	err = acct.CompleteUpdate()
	assert.NilError(t, err)
}

func testDelete(t *testing.T, acct *Context, id string) {
	err := acct.BeginUpdate()
	assert.NilError(t, err)
	defer acct.AbortUpdate()

	s, err := acct.Delete(id)
	assert.NilError(t, err)

	assert.Equal(t, s, Stats{0, 1, 0})

	err = acct.CompleteUpdate()
	assert.NilError(t, err)
}

func TestErrors(t *testing.T) {
	acct := setupTest(t)

	testRecord := Record{"", Date(time.Now()), "Pen Island", "-75.45", "Frivolities", "Tchotchkes", ""}
	s, err := acct.Insert(&testRecord)
	assert.Equal(t, s, Stats{0, 0, 0})
	assert.ErrorContains(t, err, "without beginning an update")

	testRecord.Id = "17"
	s, err = acct.Update(&testRecord)
	assert.Equal(t, s, Stats{0, 0, 0})
	assert.ErrorContains(t, err, "without beginning an update")

	s, err = acct.Delete("17")
	assert.Equal(t, s, Stats{0, 0, 0})
	assert.ErrorContains(t, err, "without beginning an update")

	err = acct.BeginUpdate()
	assert.NilError(t, err)
	defer acct.AbortUpdate()

	s, err = acct.Insert(&testRecord)
	assert.Equal(t, s, Stats{0, 0, 0})
	assert.ErrorContains(t, err, "with rowid")

	testRecord.Id = ""
	s, err = acct.Update(&testRecord)
	assert.Equal(t, s, Stats{0, 0, 0})
	assert.ErrorContains(t, err, "without rowid")
}

func TestModifications(t *testing.T) {
	acct := setupTest(t)

	//////////// Insertion

	testRecord := Record{"", Date(time.Now()), "Pen Island", "-75.45", "Frivolities", "Tchotchkes", ""}
	testInsertion(t, acct, &testRecord)

	recs := materializeQuery(t, acct, QuerySpec{})
	assert.Equal(t, len(recs), 1)
	assert.Assert(t, recs[0].Id != "")
	assertRecordsEqualNoId(t, recs[0], testRecord)

	//////////// Update

	updateRecord := Record{recs[0].Id, Date(time.Now()), "Qwik-E-Mart", "17.42", "Necessities", "Ice Cream", ""}
	testUpdate(t, acct, &updateRecord)

	recs = materializeQuery(t, acct, QuerySpec{})
	assert.Equal(t, len(recs), 1)
	assertRecordsEqual(t, recs[0], updateRecord)

	//////////// Delete

	testDelete(t, acct, updateRecord.Id)

	recs = materializeQuery(t, acct, QuerySpec{})
	assert.Equal(t, len(recs), 0)
}

func tm(s string) Date {
	d, err := ParseDate(s)
	if err != nil {
		panic("bad time constant")
	}
	return d
}

func TestQuery(t *testing.T) {
	acct := setupTest(t)

	date1 := tm("2022-11-01")
	date2 := tm("2022-11-02")
	date3 := tm("2022-11-03")

	testRecords := []Record{
		{"", date1, "Pen Island", "-75.45", "Frivolities", "Tchotchkes", ""},
		{"", date2, "Qwik-E-Mart", "17.42", "Necessities", "Ice Cream", ""},
		{"", date2, "Qwik Stop", "17.42", "Necessities", "Wine", ""},
		{"", date3, "Dewey Cheatham & Howe", "17.42", "Legalities", "Success Fees", ""}}

	for _, r := range testRecords {
		testInsertion(t, acct, &r)
	}

	recs := materializeQuery(t, acct, QuerySpec{})
	assert.Equal(t, len(recs), 4)

	necessities := "Necessities"
	recs = materializeQuery(t, acct, QuerySpec{Category: &necessities})
	assert.Equal(t, len(recs), 2)

	legalities := "Legalities"
	recs = materializeQuery(t, acct, QuerySpec{Category: &legalities})
	assert.Equal(t, len(recs), 1)

	recs = materializeQuery(t, acct, QuerySpec{Subcategory: &legalities})
	assert.Equal(t, len(recs), 0)

	icecream := "Ice Cream"
	recs = materializeQuery(t, acct, QuerySpec{Subcategory: &icecream})
	assert.Equal(t, len(recs), 1)

	recs = materializeQuery(t, acct, QuerySpec{Category: &necessities, Subcategory: &icecream})
	assert.Equal(t, len(recs), 1)

	qwik := "qwik%"
	recs = materializeQuery(t, acct, QuerySpec{DescrLike: &qwik})
	assert.Equal(t, len(recs), 2)

	recs = materializeQuery(t, acct, QuerySpec{DateFrom: &date2})
	assert.Equal(t, len(recs), 3)

	recs = materializeQuery(t, acct, QuerySpec{DateUntil: &date2})
	assert.Equal(t, len(recs), 1)
}
