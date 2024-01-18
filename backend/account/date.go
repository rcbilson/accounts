package account

import (
	"strings"
	"time"
)

type Date time.Time

func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
        return Date(t), err
}

func (c *Date) String() string {
	return time.Time(*c).Format("2006-01-02")
}

func (c *Date) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	*c = Date(t) //set result using the pointer
	return nil
}

func (c Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}
