package account

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

func parseTD(record []string) (*Record, error) {
	// 01/05/2021,SEND E-TFR CA***t5J ,413.00,,13041.20
	var err error
	var r Record
	t, err := time.Parse("01/02/2006", record[0])
	if err != nil {
                t, err = time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, err
		}
	}
	r.Date = Date(t)
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

func parseTangerine(record []string) (*Record, error) {
	// 1/1/2021,DEBIT,Recurring Internet Withdrawal to,SPOINT To 3023010786,-250
	var err error
	var r Record
	t, err := time.Parse("1/2/2006", record[0])
	if err != nil {
		return nil, err
	}
	r.Date = Date(t)
	r.Descr = record[2] + "/" + record[3]
	_, err = strconv.ParseFloat(record[4], 64)
	if err != nil {
		return nil, err
	}
	r.Amount = record[4]
	return &r, nil
}

func (acct *Context) Import(reader io.Reader) (Stats, error) {
	var stats Stats

	err := acct.BeginUpdate()
	if err != nil {
		return stats, err
	}
	defer acct.AbortUpdate()

	r := csv.NewReader(reader)

	parser := parseTD
	rowCount := 0
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
			return stats, err
		}

		err = acct.InferCategory(rec)
		if err != nil {
			return stats, err
		}

		s, err := acct.Insert(rec)
		stats.Add(s)
		if err != nil {
			return stats, err
		}
	}

	err = acct.CompleteUpdate()
	return stats, err
}
