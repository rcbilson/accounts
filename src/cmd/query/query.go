package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"

	"knilson.org/accounts/account"
)

func mustParseDate(date string) *account.Date {
	d, err := account.ParseDate(date)
	if err != nil {
		log.Fatal(err)
	}
	return &d
}

func main() {
	dateFrom := flag.String("dateFrom", "", "Earliest date for which results should be reported")
	dateUntil := flag.String("dateUntil", "", "Earliest date for which results should not be reported")
	descrLike := flag.String("descrLike", "", "Filter by matching description (SQL LIKE pattern)")
	category := flag.String("category", "", "Filter by exact match on category")
	subcategory := flag.String("subcategory", "", "Filter by exact match on subcategory")
	state := flag.String("state", "", "Filter by exact match on state")
	limit := flag.Int("limit", 0, "Maximum number of transactions to return")
	offset := flag.Int("offset", 0, "Number of transactions to skip before returning results")
	flag.Parse()

	var query account.QuerySpec
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "dateFrom" {
			query.DateFrom = mustParseDate(*dateFrom)
		} else if f.Name == "dateUntil" {
			query.DateUntil = mustParseDate(*dateUntil)
		} else if f.Name == "descrLike" {
			query.DescrLike = descrLike
		} else if f.Name == "category" {
			query.Category = category
		} else if f.Name == "subcategory" {
			query.Subcategory = subcategory
		} else if f.Name == "state" {
			query.State = state
		} else if f.Name == "limit" {
			query.Limit = limit
		} else if f.Name == "offset" {
			query.Offset = offset
		}
	})

	acct, err := account.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer acct.Close()

	ch, err := acct.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	for r := range ch {
		record := []string{r.Id, r.Date.String(), r.Descr, r.Amount, r.Category, r.Subcategory}
		writer.Write(record)
	}
}
