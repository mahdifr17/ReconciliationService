package repository

type TransactionRPInterface interface {
	ProcessTransactionCsv()
}

type BankStatementRPInterface interface {
	ProcessBankStatementCsv()
}
