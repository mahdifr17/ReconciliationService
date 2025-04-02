package usecase

import (
	"encoding/csv"
	"fmt"
	"sync"

	"github.com/mahdifr17/ReconciliationService/entity"
	"github.com/mahdifr17/ReconciliationService/repository"
	"github.com/mahdifr17/ReconciliationService/utils"
)

type ReconciliationUsecaseImpl struct {
	TransactionRP   repository.TransactionRPInterface
	BankStatementRP repository.BankStatementRPInterface
}

// Assume both internal data & bank statement are sorted by transaction time
func (uc ReconciliationUsecaseImpl) ReconcileData(
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
	var wg sync.WaitGroup
	var mu sync.Mutex
	var workerNum = len(csvBankStatements) * 2

	// using worker pool
	jobs := make(chan entity.Transaction, workerNum)

	// do work
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				// TODO: decide compare to which bank statement
				// using pattern match?
				// we should know whick provider is this bank statement
				// OR, try every bank statement
				mu.Lock()
				fmt.Println(job)
				mu.Unlock()
			}
		}()
	}

	// seed job
	for _, trx := range internalData {
		jobs <- trx
	}
	close(jobs)

	wg.Wait()
	// complete
}
