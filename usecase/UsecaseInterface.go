package usecase

import (
	"encoding/csv"
	"time"

	"github.com/mahdifr17/ReconciliationService/entity"
)

type ReconciliationUsecase interface {
	ReconcileData(csvInternalTrxData *csv.Reader, csvBankStatements []*csv.Reader, startDate, endDate time.Time) entity.ReconciliationResult
}
