package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/mahdifr17/ReconciliationService/usecase"
)

func main() {
	fmt.Println("Reconciliation Service")

	// Internal Data
	// Open file
	internalData, err := os.Open("internalTransactionData.csv")
	if err != nil {
		panic(err)
	}
	defer internalData.Close()
	// Create a new CSV reader
	readerInternalData := csv.NewReader(internalData)

	// Bank Statement
	bankStatement1, err := os.Open("bankStatementAbc.csv")
	if err != nil {
		panic(err)
	}
	defer bankStatement1.Close()
	// Create a new CSV reader
	readerBankStatement1 := csv.NewReader(bankStatement1)

	bankStatement2, err := os.Open("bankStatementDef.csv")
	if err != nil {
		panic(err)
	}
	defer bankStatement2.Close()
	// Create a new CSV reader
	readerBankStatement2 := csv.NewReader(bankStatement2)

	bankStatement3, err := os.Open("bankStatementDku.csv")
	if err != nil {
		panic(err)
	}
	defer bankStatement3.Close()
	// Create a new CSV reader
	readerBankStatement3 := csv.NewReader(bankStatement3)

	bankStatement4, err := os.Open("bankStatementZdt.csv")
	if err != nil {
		panic(err)
	}
	defer bankStatement4.Close()
	// Create a new CSV reader
	readerBankStatement4 := csv.NewReader(bankStatement4)

	var uc usecase.ReconciliationUsecase
	uc = new(usecase.ReconciliationUsecaseImpl)
	uc.ReconcileData(readerInternalData, []*csv.Reader{readerBankStatement1, readerBankStatement2, readerBankStatement3, readerBankStatement4})
}
