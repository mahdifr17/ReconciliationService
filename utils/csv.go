package utils

import (
	"encoding/csv"

	"github.com/mahdifr17/ReconciliationService/entity"
)

func LoadCsv(csvReader *csv.Reader) ([][]string, error) {
	return csvReader.ReadAll()
}

func LoadCsvInternalTrx(csvReader *csv.Reader) ([]entity.Transaction, error) {
	out := make([]entity.Transaction, 0)

	csvData, err := LoadCsv(csvReader)
	if err != nil {
		return nil, err
	}

	trx := new(entity.Transaction)
	for _, row := range csvData {
		err = trx.ReadFromCsv(row)
		if err != nil {
			// notif user this data is error
			continue
		}
		out = append(out, *trx)
	}

	return out, nil
}

func LoadCsvBankStatement(csvReader *csv.Reader) ([]entity.BankStatement, error) {
	out := make([]entity.BankStatement, 0)

	csvData, err := LoadCsv(csvReader)
	if err != nil {
		return nil, err
	}

	bankStatement := new(entity.BankStatement)
	for _, row := range csvData {
		err = bankStatement.ReadFromCsv(row)
		if err != nil {
			// notif user this data is error
			continue
		}
		out = append(out, *bankStatement)
	}

	return out, nil
}
