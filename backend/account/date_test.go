package account

import (
	"encoding/json"
	"testing"

	"gotest.tools/assert"
)

func TestDate(t *testing.T) {
	_, err := ParseDate("2021-08-12T00:00:00Z")
	assert.Assert(t, err != nil)

	d, err := ParseDate("2021-08-12")
	assert.Equal(t, "2021-08-12", d.String())

	type testJson struct {
		Date Date
	}
	strct := testJson{d}
	js, err := json.Marshal(strct)
	assert.NilError(t, err)
        var newstrct testJson
        err = json.Unmarshal(js, &newstrct)
        assert.NilError(t, err)
        assert.Equal(t, d, newstrct.Date)
}
