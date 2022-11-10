package main

import (
	"encoding/csv"
	"log"
	"os"

	"knilson.org/accounts/account"
)

func main() {
	acct, err := account.Open()
	if err != nil {
		log.Fatal(err)
	}
        defer acct.Close()

	where := os.Args[1]
	ch, err := acct.Query(where)
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	for r := range ch {
		record := []string{r.Id, r.Date.Format("2006-01-02"), r.Descr, r.Amount, r.Category, r.Subcategory}
		writer.Write(record)
	}
}
