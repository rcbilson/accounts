package account

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestParseTD(t *testing.T) {
	input := [][]string{
		{"03/02/2020", "E-TRANSFER CA***EGh ", "", "317.77", "11829.70"},
		{"10/21/2022", "WRDSB PAYROLL    PAY", "", "2590.51", "28614.06"},
		{"10/17/2022", "UT172 TFR-TO C/C    ", "4238.11", "", "26591.95"},
	}
	expected := []Record{
		{"", tm("2020-03-02"), "E-TRANSFER CA***EGh ", "317.77", "", "", ""},
		{"", tm("2022-10-21"), "WRDSB PAYROLL    PAY", "2590.51", "", "", ""},
		{"", tm("2022-10-17"), "UT172 TFR-TO C/C    ", "-4238.11", "", "", ""},
	}
	for i := range input {
		r, err := parseTD(input[i])
		assert.NilError(t, err)
		assert.Equal(t, *r, expected[i])
	}
}

func TestParseTDError(t *testing.T) {
	input := [][]string{
		{"14/02/2020", "E-TRANSFER CA***EGh ", "", "317.77", "11829.70"},
		{"2022-10-21", "WRDSB PAYROLL    PAY", "", "2590.51", "28614.06"},
		{"10/21/2022", "WRDSB PAYROLL    PAY", "", "seventeen", "28614.06"}}
	for i := range input {
		_, err := parseTD(input[i])
		assert.Assert(t, err != nil, fmt.Sprintf("index %v", i))
	}
}

func TestParseTangerine(t *testing.T) {
	input := [][]string{
		{"9/1/2022", "OTHER", "EFT Deposit from THE TORONTO-DOM", "COFA", "200"},
		{"9/14/2022", "ATM", "EFT Tangerine Credit Card Paymen", "To TANGERINE CCRD", "-215.75"},
	}
	expected := []Record{
		{"", tm("2022-09-01"), "EFT Deposit from THE TORONTO-DOM/COFA", "200", "", "", ""},
		{"", tm("2022-09-14"), "EFT Tangerine Credit Card Paymen/To TANGERINE CCRD", "-215.75", "", "", ""},
	}
	for i := range input {
		r, err := parseTangerine(input[i])
		assert.NilError(t, err)
		assert.Equal(t, *r, expected[i])
	}
}

func TestParseTangerineError(t *testing.T) {
	input := [][]string{
		{"9/1/2022", "OTHER", "EFT Deposit from THE TORONTO-DOM", "COFA", "seventeen"},
		{"14/9/2022", "ATM", "EFT Tangerine Credit Card Paymen", "To TANGERINE CCRD", "-215.75"}}
	for i := range input {
		_, err := parseTD(input[i])
		assert.Assert(t, err != nil, fmt.Sprintf("index %v", i))
	}
}
