package usecase

import (
	"encoding/csv"
)

type ReconciliationUsecase interface {
	ReconcileData(csvInternalTrxData *csv.Reader, csvBankStatements []*csv.Reader)
}
