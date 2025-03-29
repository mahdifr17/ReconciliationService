package entity

const (
	TrxTypeDebit  TransactionType = iota // enum idx 0
	TrxTypeCredit                        // enum idx 1
)

type TransactionType int

func (t TransactionType) String() string {
	return [...]string{"Debit", "Credit"}[t]
}

func (t TransactionType) EnumIndex() int {
	return int(t)
}

type TransactionType2 string

const (
	TrxTypeDebit  TransactionType2 = "DEBIT"
	TrxTypeCredit TransactionType2 = "CREDIT"
)

func IsValidTrxType(input string) (TransactionType2, error) {
	
}