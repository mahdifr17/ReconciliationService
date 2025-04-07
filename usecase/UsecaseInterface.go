package usecase

import (
	"encoding/csv"

	"github.com/mahdifr17/ReconciliationService/entity"
)

type ReconciliationUsecase interface {
	ReconcileData(csvInternalTrxData *csv.Reader, csvBankStatements []*csv.Reader) []entity.ReconciliationResult
}
