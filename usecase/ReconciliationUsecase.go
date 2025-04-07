package usecase

import (
	"encoding/csv"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/mahdifr17/ReconciliationService/entity"
	"github.com/mahdifr17/ReconciliationService/utils"
)

const (
	errCaseNotFoundOnInternal      = "Not Found on Internal Transaction"
	errCaseNotFoundOnBankStatement = "Not Found on Bank Statement"
	errCaseAmountMismatch          = "Amount Mismatch"
)

type ReconciliationUsecaseImpl struct {
}

// Assume both internal data & bank statement are sorted by transaction time
// Reconcile
// using worker pool
func (uc ReconciliationUsecaseImpl) ReconcileData(
	csvInternalTrxData *csv.Reader, csvBankStatements []*csv.Reader,
) []entity.ReconciliationResult {
	fmt.Println("Start Reconcile")

	var (
		out                  = make([]entity.ReconciliationResult, 0)
		timeout              = time.After(1 * time.Minute)
		wg                   sync.WaitGroup
		mu                   sync.Mutex
		csvBuffer            = 1000
		workerPerBs          = 3
		mapInternalTrx       = make(map[string]*entity.Transaction)
		streamTrxInternal    = make(chan entity.Transaction, csvBuffer)
		streamBankStatements = make([]chan entity.BankStatement, 0)
		chanFinish           = make(chan bool)
	)

	// seed job
	// read data internal transaction
	go utils.LoadCsvInternalTrx(csvInternalTrxData, streamTrxInternal)
	for trxInternal := range streamTrxInternal {
		mapInternalTrx[trxInternal.TrxId] = &trxInternal
	}
	// read data bank statement
	for _, csvBankStatement := range csvBankStatements {
		streamBankStatement := make(chan entity.BankStatement, csvBuffer)
		streamBankStatements = append(streamBankStatements, streamBankStatement)
		go utils.LoadCsvBankStatement(csvBankStatement, streamBankStatement)
	}

	// do work
	for _, streamBankStatement := range streamBankStatements {
		for i := 0; i < workerPerBs; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for bankStatement := range streamBankStatement {
					mu.Lock()
					result := uc.compareData(bankStatement, mapInternalTrx)
					if result.Remark != "" {
						out = append(out, result)
					}
					mu.Unlock()
				}
			}()
		}
	}

	go func() {
		wg.Wait()
		chanFinish <- true
	}()

	select {
	case <-timeout:
		fmt.Println("Reconcile Timeout!")
	case <-chanFinish:
		// complete
		// check is there any more internal trx, if this happen, then trx not found on bank statement
		result := uc.checkRemainingTrxInternal(mapInternalTrx)
		if len(result) != 0 {
			out = append(out, result...)
		}
		fmt.Println("Finish Reconcile")
	}
	return out
}

// Cases
// exist on bank statement - not found on internal
// different amount
func (uc ReconciliationUsecaseImpl) compareData(
	bankStatement entity.BankStatement, mapInternalTrx map[string]*entity.Transaction,
) entity.ReconciliationResult {
	trx, exist := mapInternalTrx[bankStatement.UniqueIdentifier]
	if !exist {
		return entity.ReconciliationResult{TrxId: bankStatement.UniqueIdentifier, Remark: errCaseNotFoundOnInternal}
	}

	compareAmount := bankStatement.Amount
	if trx.Type == entity.TrxTypeDebit {
		compareAmount = math.Abs(bankStatement.Amount)
	}

	if compareAmount != trx.Amount {
		trx.Remark = errCaseAmountMismatch
		return entity.ReconciliationResult{TrxId: bankStatement.UniqueIdentifier, Remark: errCaseAmountMismatch}
	}

	delete(mapInternalTrx, bankStatement.UniqueIdentifier)
	return entity.ReconciliationResult{}
}

// Cases
// not found on bank statement - exist on internal
func (uc ReconciliationUsecaseImpl) checkRemainingTrxInternal(
	mapInternalTrx map[string]*entity.Transaction,
) (out []entity.ReconciliationResult) {
	out = make([]entity.ReconciliationResult, 0)
	for _, internalTrx := range mapInternalTrx {
		if internalTrx.Remark != errCaseAmountMismatch {
			out = append(out, entity.ReconciliationResult{TrxId: internalTrx.TrxId, Remark: errCaseNotFoundOnBankStatement})
		}
	}
	return out
}
