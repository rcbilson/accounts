package main

import (
	"log"
	"os"

	"knilson.org/accounts/account"
)

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
	check(err)
	defer acct.Close()

	var stats account.Stats
	for _, csvFile := range os.Args[1:] {
		f, err := os.Open(csvFile)
		check(err)

		s, err := acct.Import(f)
		check(err)
		stats.Add(s)
	}
	log.Printf("done %v", stats)
}
