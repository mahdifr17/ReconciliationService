package entity

import (
	"strconv"
	"time"
)

type Transaction struct {
	TrxId           string
	Amount          float64
	Type            TransactionType
	TransactionTime time.Time // datetime format
}

func (t *Transaction) ReadFromCsv(input []string) (err error) {
	// input must be trx_id, amount, type, transaction_time
	if len(input) != 4 {
		return
	}

	t.TrxId = input[0]
	t.Amount, err = strconv.ParseFloat(input[1], 64)
	if err != nil {
		return err
	}
	switch input[2] {
	case TrxTypeDebit.String():
		t.Type = TrxTypeDebit
	case TrxTypeCredit.String():
		t.Type = TrxTypeCredit
	default:
		return ErrorInvalidTransactionType{input[2]}
	}
	t.TransactionTime, err = time.Parse(time.DateTime, input[3])
	if err != nil {
		return err
	}
	return nil
}
