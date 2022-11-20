package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"time"

	"knilson.org/accounts/account"
)

func mustParseDate(date string) *time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Fatal(err)
	}
	return &t
}

func main() {
	dateFrom := flag.String("dateFrom", "", "Earliest date for which results should be reported")
	dateUntil := flag.String("dateUntil", "", "Earliest date for which results should not be reported")
	descrLike := flag.String("descrLike", "", "Filter by matching description (SQL LIKE pattern)")
	category := flag.String("category", "", "Filter by exact match on category")
	subcategory := flag.String("subcategory", "", "Filter by exact match on subcategory")
	state := flag.String("state", "", "Filter by exact match on state")
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
		record := []string{r.Id, r.Date.Format("2006-01-02"), r.Descr, r.Amount, r.Category, r.Subcategory}
		writer.Write(record)
	}
}
