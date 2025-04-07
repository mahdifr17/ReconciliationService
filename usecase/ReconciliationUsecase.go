package usecase

import (
	"encoding/csv"
	"fmt"
	"math"
	"strconv"
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
	csvInternalTrxData *csv.Reader, csvBankStatements []*csv.Reader, startDate, endDate time.Time,
) entity.ReconciliationResult {
	fmt.Println("Start Reconcile")

	var (
		out                  = new(entity.ReconciliationResult)
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
	out.ListMissingTransactionInternal = make(map[string][]entity.CompareResult)

	// seed job
	// read data internal transaction
	go utils.LoadCsvInternalTrx(csvInternalTrxData, streamTrxInternal)
	for trxInternal := range streamTrxInternal {
		if isInRangeDate(trxInternal.TransactionTime, startDate, endDate) {
			mapInternalTrx[trxInternal.TrxId] = &trxInternal
		}
	}
	// read data bank statement
	for _, csvBankStatement := range csvBankStatements {
		streamBankStatement := make(chan entity.BankStatement, csvBuffer)
		streamBankStatements = append(streamBankStatements, streamBankStatement)
		go utils.LoadCsvBankStatement(csvBankStatement, streamBankStatement)
	}

	// do work
	for idx, streamBankStatement := range streamBankStatements {
		for i := 0; i < workerPerBs; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for bankStatement := range streamBankStatement {
					if !isInRangeDate(bankStatement.Date.Add(5*time.Minute), startDate, endDate) {
						continue
					}

					mu.Lock()
					result := uc.compareData(bankStatement, mapInternalTrx)
					out.TotalTransactionProcessed++
					if result.Remark == "" {
						out.TotalMatchTransaction++
					} else if result.Remark == errCaseAmountMismatch {
						diff := result.InternalTransaction.Amount - result.BankStatement.Amount
						out.TotalDiscrepancies += math.Abs(diff)
					} else {
						sIdx := strconv.Itoa(idx)
						if _, exist := out.ListMissingTransactionInternal[sIdx]; !exist {
							out.ListMissingTransactionInternal[sIdx] = make([]entity.CompareResult, 0)
						}
						out.ListMissingTransactionInternal[sIdx] = append(out.ListMissingTransactionInternal[sIdx], result)
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
			out.ListMissingTransactionBank = result
		}
		fmt.Println("Finish Reconcile")
	}
	return *out
}

func isInRangeDate(check, startDate, endDate time.Time) bool {
	if check.After(startDate) && check.Before(endDate) {
		return true
	}
	return false
}

// Cases
// exist on bank statement - not found on internal
// different amount
func (uc ReconciliationUsecaseImpl) compareData(
	bankStatement entity.BankStatement, mapInternalTrx map[string]*entity.Transaction,
) entity.CompareResult {
	trx, exist := mapInternalTrx[bankStatement.UniqueIdentifier]
	if !exist {
		return entity.CompareResult{BankStatement: bankStatement, Remark: errCaseNotFoundOnInternal}
	}

	compareAmount := bankStatement.Amount
	if trx.Type == entity.TrxTypeDebit {
		compareAmount = math.Abs(bankStatement.Amount)
	}

	if compareAmount != trx.Amount {
		trx.Remark = errCaseAmountMismatch
		return entity.CompareResult{InternalTransaction: *trx, BankStatement: bankStatement, Remark: errCaseAmountMismatch}
	}

	delete(mapInternalTrx, bankStatement.UniqueIdentifier)
	return entity.CompareResult{}
}

// Cases
// not found on bank statement - exist on internal
func (uc ReconciliationUsecaseImpl) checkRemainingTrxInternal(
	mapInternalTrx map[string]*entity.Transaction,
) (out []entity.CompareResult) {
	out = make([]entity.CompareResult, 0)
	for _, internalTrx := range mapInternalTrx {
		if internalTrx.Remark != errCaseAmountMismatch {
			out = append(out, entity.CompareResult{InternalTransaction: *internalTrx, Remark: errCaseNotFoundOnBankStatement})
		}
	}
	return out
}
