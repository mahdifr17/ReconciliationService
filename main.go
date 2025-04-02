package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/mahdifr17/ReconciliationService/usecase"
)

func main2() {
	fmt.Println("Reconciliation Service")

	// Internal Data
	// Open file
	internalData, err := os.Open("internalTransaction.csv")
	if err != nil {
		panic(err)
	}
	defer internalData.Close()

	// Create a new CSV reader
	readerInternalData := csv.NewReader(internalData)

	// Bank Statement
	bankStatement, err := os.Open("bankStatement1.csv")
	if err != nil {
		panic(err)
	}
	defer bankStatement.Close()

	// Create a new CSV reader
	readerBankStatement := csv.NewReader(bankStatement)

	var uc usecase.ReconciliationUsecase
	uc = new(usecase.ReconciliationUsecaseImpl)
	uc.ReconcileData(readerInternalData, []*csv.Reader{readerBankStatement})
}
