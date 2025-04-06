package usecase

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/mahdifr17/ReconciliationService/entity"
	"github.com/mahdifr17/ReconciliationService/utils"
)

type SimpleReconciliationUsecaseImpl struct {
}

// Assume both internal data & bank statement are sorted by transaction time
func (sr SimpleReconciliationUsecaseImpl) ReconcileData(
	csvInternalTrxData *csv.Reader, csvBankStatements []*csv.Reader,
) {
	// Load the data
	internalData, err := utils.LoadCsvInternalTrx(csvInternalTrxData)
	if err != nil {
		// TODO: err handling
	}

	bankStatements := make([][]entity.BankStatement, 0)
	for _, csvBankStatement := range csvBankStatements {
		bankStatement, err := utils.LoadCsvBankStatement(csvBankStatement)
		if err != nil {
			// TODO: err handling
			continue
		}
		bankStatements = append(bankStatements, bankStatement)
	}

	// Reconcile
	for _, internalTrx := range internalData {
		for bankStatementsIdx, bankStatement := range bankStatements {
			if len(bankStatement) > 0 {
				if strings.EqualFold(internalTrx.TrxId, bankStatement[0].UniqueIdentifier) {
					// found
					bankStatements[bankStatementsIdx] = bankStatement[1:] // remove from bank statement
					fmt.Printf("Found: %v, %v\n", internalTrx.TrxId, bankStatementsIdx)
					break
				}
			}
			// else continue check other bank statement
			continue
		}
	}

	fmt.Println("Finish")
}
