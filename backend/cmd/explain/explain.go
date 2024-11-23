package main

import (
	"flag"
	"fmt"
	"log"

	"knilson.org/accounts/account"
)

func check(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	acct, err := account.Open()
	check(err)
	defer acct.Close()

	descr := flag.String("descr", "", "description from bank")
	amount := flag.String("amount", "", "transaction amount")
	flag.Parse()
        fmt.Println("Amount", *amount, "Description", *descr)
	rec := account.Record{Descr: *descr, Amount: *amount}

	err, explanation := acct.InferCategoryWithExplanation(&rec)
	check(err)

	fmt.Println("Inferred category", rec.Category, "/", rec.Subcategory)
	fmt.Println("explanation:", explanation)
}
