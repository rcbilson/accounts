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

func materializeQuery(t *testing.T, acct *Context, query string) []Record {
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
	assert.Equal(t, left.Date.Format("2006-01-02"), right.Date.Format("2006-01-02"))
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

	testRecord := Record{"", time.Now(), "Pen Island", "-75.45", "Frivolities", "Tchotchkes"}
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

	testRecord := Record{"", time.Now(), "Pen Island", "-75.45", "Frivolities", "Tchotchkes"}
	testInsertion(t, acct, &testRecord)

	recs := materializeQuery(t, acct, "")
	assert.Equal(t, len(recs), 1)
	assert.Assert(t, recs[0].Id != "")
	assertRecordsEqualNoId(t, recs[0], testRecord)

	//////////// Update

	updateRecord := Record{recs[0].Id, time.Now(), "Qwik-E-Mart", "17.42", "Necessities", "Ice Cream"}
	testUpdate(t, acct, &updateRecord)

	recs = materializeQuery(t, acct, "")
	assert.Equal(t, len(recs), 1)
	assertRecordsEqual(t, recs[0], updateRecord)

	//////////// Delete

	testDelete(t, acct, updateRecord.Id)

	recs = materializeQuery(t, acct, "")
	assert.Equal(t, len(recs), 0)
}
