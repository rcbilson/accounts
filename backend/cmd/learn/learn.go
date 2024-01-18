package main

import (
	"log"

	"knilson.org/accounts/account"
)

type Specification struct {
	DbFile string
}

func check(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func checkRow(err error, rowCount int) {
	if err != nil {
		log.Fatalf("at line %d: %v", rowCount, err.Error())
	}
}

func main() {
	acct, err := account.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer acct.Close()

	acct.ResetLearning()

	ch, err := acct.Query(account.QuerySpec{})
	check(err)

	err = acct.BeginUpdate()
	check(err)
	defer acct.AbortUpdate()

	rowCount := 0
	var stats account.Stats
	for r := range ch {
		rowCount++
		s, err := acct.UpdateLearning(r)
		checkRow(err, rowCount)
		stats.Add(s)
	}
	err = acct.CompleteUpdate()
	check(err)

	log.Println(rowCount, "rows processed", stats)
}
