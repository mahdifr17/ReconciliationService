package usecase

import (
	"github.com/mahdifr17/ReconciliationService/repository"
)

type ReconciliationUsecaseImpl struct {
	TransactionRP   repository.TransactionRPInterface
	BankStatementRP repository.BankStatementRPInterface
}

func (uc ReconciliationUsecaseImpl) ReconcileData(
	internalTrxData multipart.File, bankStatement []multipart.File,
	) {
		
}
