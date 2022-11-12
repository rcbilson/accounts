package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

func parseTD(record []string) (*account.Record, error) {
	// 01/05/2021,SEND E-TFR CA***t5J ,413.00,,13041.20
	var err error
	var r account.Record
	r.Date, err = time.Parse("01/02/2006", record[0])
	if err != nil {
		return nil, err
	}
	r.Descr = record[1]
	debit, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		debit = 0
	}
	credit, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		credit = 0
	}
	if credit == 0 && debit == 0 {
		return nil, errors.New("No valid credit/debit")
	}
	r.Amount = fmt.Sprintf("%.2f", credit-debit)
	return &r, nil
}

func parseTangerine(record []string) (*account.Record, error) {
	// 1/1/2021,DEBIT,Recurring Internet Withdrawal to,SPOINT To 3023010786,-250
	var err error
	var r account.Record
	r.Date, err = time.Parse("1/2/2006", record[0])
	if err != nil {
		return nil, err
	}
	r.Descr = record[2] + "/" + record[3]
	_, err = strconv.ParseFloat(record[4], 64)
	if err != nil {
		return nil, err
	}
	r.Amount = record[4]
	return &r, nil
}

func main() {
	acct, err := account.Open()
	check(err)
	defer acct.Close()

	err = acct.BeginUpdate()
	check(err)
	defer acct.AbortUpdate()

	for _, csvFile := range os.Args[1:] {
		f, err := os.Open(csvFile)
		check(err)

		r := csv.NewReader(f)

		parser := parseTD
		rowCount := 0
		var stats account.Stats
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			rowCount++
			if record[0] == "Date" {
				parser = parseTangerine
				continue
			}
			rec, err := parser(record)
			if err != nil {
				log.Fatalf("at %s row %d: %v", csvFile, rowCount, err)
			}

			err = acct.InferCategory(rec)
			checkRow(err, rowCount)

			s, err := acct.Insert(rec)
			checkRow(err, rowCount)
			stats.Add(s)
		}

		log.Printf("%s: %v rows processed %v", csvFile, rowCount, stats)
	}
	err = acct.CompleteUpdate()
	check(err)
}
