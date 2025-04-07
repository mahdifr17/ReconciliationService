package usecase

import (
	"encoding/csv"
	"fmt"
	"math"
	"sync"
	"time"

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
	fmt.Println("Start Reconcile")

	// Reconcile
	// using worker pool
	var (
		wg                   sync.WaitGroup
		mu                   sync.Mutex
		cond                 = sync.NewCond(&mu)
		csvBuffer            = 1000
		workerNumPerBs       = 3
		mapTrxInternalBuffer = 20
		mapInternalTrx       = make(map[string]entity.Transaction)
		streamTrxInternal    = make(chan entity.Transaction, csvBuffer)
		streamBankStatements = make([]chan entity.BankStatement, 0)
	)

	// seed job
	// read data internal transaction
	go utils.LoadCsvInternalTrx(csvInternalTrxData, streamTrxInternal)
	go func() {
		for trxInternal := range streamTrxInternal {
			cond.L.Lock()
			mapInternalTrx[trxInternal.TrxId] = trxInternal
			if len(mapInternalTrx) > len(csvBankStatements)*workerNumPerBs*mapTrxInternalBuffer {
				cond.Broadcast()
			}
			cond.L.Unlock()
		}
		cond.Broadcast()
	}()
	time.Sleep(1 * time.Second)
	// read data bank statement
	for _, csvBankStatement := range csvBankStatements {
		streamBankStatement := make(chan entity.BankStatement, csvBuffer)
		streamBankStatements = append(streamBankStatements, streamBankStatement)
		go utils.LoadCsvBankStatement(csvBankStatement, streamBankStatement)
	}

	// do work
	for _, streamBankStatement := range streamBankStatements {
		for i := 0; i < workerNumPerBs; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for bankStatement := range streamBankStatement {
					cond.L.Lock() // for locking mapInternalTrx
					for len(mapInternalTrx) < mapTrxInternalBuffer && len(streamTrxInternal) > 0 {
						cond.Wait()
					}
					uc.compareData(bankStatement, mapInternalTrx)
					cond.L.Unlock()
				}
			}()
		}
	}

	wg.Wait()
	// complete

	// check is there any more internal trx, if this happen, then trx not found on bank statement
	uc.checkRemainingTrxInternal(mapInternalTrx)
	fmt.Println("Finish Reconcile")
}

// Cases
// exist on bank statement - not found on internal
// not found on bank statement - exist on internal
// different amount
func (uc ReconciliationUsecaseImpl) compareData(
	bankStatement entity.BankStatement, mapInternalTrx map[string]entity.Transaction,
) bool {
	trx, exist := mapInternalTrx[bankStatement.UniqueIdentifier]
	if !exist {
		fmt.Println("Not Found on internal transaction", bankStatement)
		return false
	}

	compareAmount := bankStatement.Amount
	if trx.Type == entity.TrxTypeDebit {
		compareAmount = math.Abs(bankStatement.Amount)
	}

	if compareAmount != trx.Amount {
		fmt.Println("Invalid amount", bankStatement, trx)
		return false
	}

	delete(mapInternalTrx, bankStatement.UniqueIdentifier)
	return true
}

// Cases
// exist on bank statement - not found on internal
// not found on bank statement - exist on internal
// different amount
func (uc ReconciliationUsecaseImpl) checkRemainingTrxInternal(
	mapInternalTrx map[string]entity.Transaction,
) {
	for _, internalTrx := range mapInternalTrx {
		fmt.Println("Not Found on bank statement", internalTrx)
	}
}
