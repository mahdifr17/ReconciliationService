package entity

type CompareResult struct {
	InternalTransaction Transaction
	BankStatement       BankStatement
	Remark              string
}

type ReconciliationResult struct {
	TotalTransactionProcessed      int64 // calculate from internal trx
	TotalMatchTransaction          int64
	ListMissingTransactionInternal map[string][]CompareResult // list of transaction that missing on internal, mapping with idx provider
	ListMissingTransactionBank     []CompareResult
	TotalDiscrepancies             float64
}
