package entity

import (
	"strconv"
	"time"
)

type BankStatement struct {
	UniqueIdentifier string
	Amount           float64   // negative value means Debit
	Date             time.Time // date
}

func (bs *BankStatement) ReadFromCsv(input []string) (err error) {
	// input must be unique_identifier, amount, date
	if len(input) != 3 {
		return
	}

	bs.UniqueIdentifier = input[0]
	bs.Amount, err = strconv.ParseFloat(input[1], 64)
	if err != nil {
		return err
	}
	bs.Date, err = time.Parse(time.DateOnly, input[2])
	if err != nil {
		return err
	}
	return nil
}
