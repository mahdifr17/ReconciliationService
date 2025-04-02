package entity

import "fmt"

const (
	TrxTypeDebit  TransactionType = iota // enum idx 0
	TrxTypeCredit                        // enum idx 1
)

type TransactionType int

func (t TransactionType) String() string {
	return [...]string{"DEBIT", "CREDIT"}[t]
}

func (t TransactionType) EnumIndex() int {
	return int(t)
}

// Define error invalid transaction type
type ErrorInvalidTransactionType struct {
	TransactionType string
}

func (e ErrorInvalidTransactionType) Error() string {
	return fmt.Sprintf("invalid transaction type: %s", e.TransactionType)
}
