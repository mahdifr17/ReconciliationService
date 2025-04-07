package utils

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/mahdifr17/ReconciliationService/entity"
)

func LoadCsvInternalTrx(csvReader *csv.Reader, streamTrx chan<- entity.Transaction) {
	_, errHeader := csvReader.Read() // header
	if errHeader == io.EOF {
		close(streamTrx)
		return
	}
	for {
		row, errRead := csvReader.Read()
		switch errRead {
		case io.EOF:
			close(streamTrx)
			return
		case nil:
			// continue
		default: // uncatch error
			fmt.Println("error read interal csv", errRead.Error())
			close(streamTrx)
			panic(errRead)
		}

		trx := new(entity.Transaction)
		errData := trx.ReadFromCsv(row)
		if errData != nil {
			fmt.Println("invalid interal data", row, errData.Error())
			continue
		}

		streamTrx <- *trx
	}
}

func LoadCsvBankStatement(csvReader *csv.Reader, streamBankStatement chan<- entity.BankStatement) {
	csvReader.Read() // header
	for {
		row, errRead := csvReader.Read()
		switch errRead {
		case io.EOF:
			close(streamBankStatement)
			return
		case nil:
			// continue
		default: // uncatch error
			fmt.Println("error read bank statement csv", errRead.Error())
			close(streamBankStatement)
			panic(errRead)
		}

		bankStatement := new(entity.BankStatement)
		errData := bankStatement.ReadFromCsv(row)
		if errData != nil {
			fmt.Println("invalid bank statement data", row, errData.Error())
			continue
		}

		streamBankStatement <- *bankStatement
	}
}
