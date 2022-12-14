package account

import (
	"regexp"
	"testing"
	"time"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func TestBuildQuery(t *testing.T) {
	query, params := buildQuery(QuerySpec{}, baseQuery, orderBy)
	assert.Assert(t, cmp.Contains(query, baseQuery))
	assert.Equal(t, 0, len(params))

	xxx := "xxx"
	query, params = buildQuery(QuerySpec{Category: &xxx}, baseQuery, orderBy)
	assert.Assert(t, cmp.Contains(query, " WHERE category = ?"))
	assert.Equal(t, params[0], xxx)

	yyy := "yyy"
	query, params = buildQuery(QuerySpec{Category: &xxx, Subcategory: &yyy}, baseQuery, orderBy)
	assert.Assert(t, cmp.Regexp(regexp.MustCompile(baseQuery+" WHERE .* AND "), query))
	assert.Assert(t, cmp.Contains(query, "category = ?"))
	assert.Assert(t, cmp.Contains(query, "subcategory = ?"))
	assert.Assert(t, cmp.Contains(params, xxx))
	assert.Assert(t, cmp.Contains(params, yyy))

	date := Date(time.Now())
	query, params = buildQuery(QuerySpec{DateFrom: &date}, baseQuery, orderBy)
	assert.Assert(t, cmp.Contains(query, " WHERE date >= ?"))
	assert.Equal(t, params[0], date.String())

	query, params = buildQuery(QuerySpec{DateUntil: &date}, baseQuery, orderBy)
	assert.Assert(t, cmp.Contains(query, " WHERE date < ?"))
	assert.Equal(t, params[0], date.String())

	query, params = buildQuery(QuerySpec{DescrLike: &xxx}, baseQuery, orderBy)
	assert.Assert(t, cmp.Contains(query, " WHERE descr like ?"))
	assert.Equal(t, params[0], xxx)

	query, params = buildQuery(QuerySpec{State: &xxx}, baseQuery, orderBy)
	assert.Assert(t, cmp.Contains(query, " WHERE state = ?"))
	assert.Equal(t, params[0], xxx)

        limit := 1
        offset := 2
	query, params = buildQuery(QuerySpec{Limit: &limit, Offset: &offset}, baseQuery, orderBy)
	assert.Assert(t, cmp.Contains(query, " LIMIT ? OFFSET ?"))
	assert.Equal(t, params[0], limit)
	assert.Equal(t, params[1], offset)
}
