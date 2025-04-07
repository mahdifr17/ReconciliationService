package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/mahdifr17/ReconciliationService/usecase"
)

func main() {
	fmt.Println("Reconciliation Service")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter internal transaction csv location: ")
	inputInternalDataLoc, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	inputInternalDataLoc = strings.TrimSpace(inputInternalDataLoc)

	// Internal Data
	// Open file
	internalData, err := os.Open(inputInternalDataLoc)
	if err != nil {
		panic(err)
	}
	defer internalData.Close()
	// Create a new CSV reader
	readerInternalData := csv.NewReader(internalData)

	// Bank Statement
	fmt.Print("Enter multiple bank statement csv location (separate with ';'): ")
	inputBankStatementsLoc, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	inputBankStatementsLoc = strings.TrimSpace(inputBankStatementsLoc)
	bankStatementsLoc := strings.Split(inputBankStatementsLoc, ";")

	listReaderBankStatement := make([]*csv.Reader, 0)
	for _, bankStatementLoc := range bankStatementsLoc {
		bankStatement, err := os.Open(bankStatementLoc)
		if err != nil {
			panic(err)
		}
		defer bankStatement.Close()
		// Create a new CSV reader
		readerBankStatement := csv.NewReader(bankStatement)
		listReaderBankStatement = append(listReaderBankStatement, readerBankStatement)
	}

	var uc usecase.ReconciliationUsecase
	uc = new(usecase.ReconciliationUsecaseImpl)
	result := uc.ReconcileData(readerInternalData, listReaderBankStatement)
	for _, v := range result {
		fmt.Printf("%s Remark: %s\n", v.TrxId, v.Remark)
	}
}
