package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

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

	// Start date
	fmt.Print("Start reconciliation date (yyyy-mm-dd): ")
	inputStartDate, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	inputStartDate = strings.TrimSpace(inputStartDate)
	startDate, err := time.Parse(time.DateOnly, inputStartDate)
	if err != nil {
		panic(err)
	}

	// End date
	fmt.Print("End reconciliation date (yyyy-mm-dd): ")
	inputEndDate, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	inputEndDate = strings.TrimSpace(inputEndDate)
	endDate, err := time.Parse(time.DateOnly, inputEndDate)
	if err != nil {
		panic(err)
	}
	endDate = endDate.Add(86399 * time.Second) // 23:59:59
	if startDate.After(endDate) {
		panic("invalid start-end date")
	}

	var uc usecase.ReconciliationUsecase
	uc = new(usecase.ReconciliationUsecaseImpl)
	result := uc.ReconcileData(readerInternalData, listReaderBankStatement, startDate, endDate)
	fmt.Printf("Total Transaction Processed: %d\n", result.TotalTransactionProcessed)
	fmt.Printf("Total Match Transaction: %d\n", result.TotalMatchTransaction)
	fmt.Printf("Detail Unmatch Transaction:\n")
	fmt.Printf("Missing in Bank Statement: %+v\n", result.ListMissingTransactionBank)
	fmt.Printf("Missing in Internal Transaction: %+v\n", result.ListMissingTransactionInternal)
	fmt.Printf("Total Discrepancies: %v\n", result.TotalDiscrepancies)
}
