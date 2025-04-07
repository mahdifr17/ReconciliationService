package usecase

import (
	"bytes"
	"encoding/csv"
	"reflect"
	"testing"

	"github.com/mahdifr17/ReconciliationService/entity"
)

func TestReconciliationUsecaseImpleeeee_ReconcileData(t *testing.T) {
	type args struct {
		csvInternalTrxData *csv.Reader
		csvBankStatements  []*csv.Reader
	}
	tests := []struct {
		name string
		uc   ReconciliationUsecase
		args args
		want []entity.ReconciliationResult
	}{
		{
			name: "valid trx",
			uc:   new(ReconciliationUsecaseImpl),
			args: args{getMockInternalTrxDataTcValid(), getMockBankStatementDataTcValid()},
			want: []entity.ReconciliationResult{},
		},
		{
			name: "invalid amount",
			uc:   new(ReconciliationUsecaseImpl),
			args: args{getMockInternalTrxDataTcValid(), getMockBankStatementDataTcInvalidAmount()},
			want: []entity.ReconciliationResult{{TrxId: "DKU202501011536751", Remark: errCaseAmountMismatch}},
		},
		{
			name: "not found on internal",
			uc:   new(ReconciliationUsecaseImpl),
			args: args{getMockInternalTrxDataTcNotFound(), getMockBankStatementDataTcValid()},
			want: []entity.ReconciliationResult{{TrxId: "DKU202501011536751", Remark: errCaseNotFoundOnInternal}},
		},
		{
			name: "not found on bank statement",
			uc:   new(ReconciliationUsecaseImpl),
			args: args{getMockInternalTrxDataTcValid(), getMockBankStatementDataTcNotFound()},
			want: []entity.ReconciliationResult{{TrxId: "DKU202501011726759", Remark: errCaseNotFoundOnBankStatement}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uc.ReconcileData(tt.args.csvInternalTrxData, tt.args.csvBankStatements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconciliationUsecaseImpl.ReconcileData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getMockInternalTrxDataTcValid() *csv.Reader {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	data := [][]string{
		{"trxId", "amount", "type", "transactionTime"},
		{"ABC202501010396818", "33160.38", "DEBIT", "2025-01-01 00:00:06"},
		{"DKU202501011536751", "84844.94", "CREDIT", "2025-01-01 00:00:08"},
		{"DKU202501011726759", "92854.13", "CREDIT", "2025-01-01 00:00:33"},
	}
	w.WriteAll(data)
	w.Flush()

	return csv.NewReader(buf)
}

func getMockBankStatementDataTcValid() []*csv.Reader {
	out := make([]*csv.Reader, 0)

	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	data := [][]string{
		{"unique_identifier", "amount", "date"},
		{"ABC202501010396818", "-33160.38", "2025-01-01"},
	}
	w.WriteAll(data)
	w.Flush()
	out = append(out, csv.NewReader(buf))

	buf2 := &bytes.Buffer{}
	w2 := csv.NewWriter(buf2)
	data2 := [][]string{
		{"unique_identifier", "amount", "date"},
		{"DKU202501011536751", "84844.94", "2025-01-01"},
		{"DKU202501011726759", "92854.13", "2025-01-01"},
	}
	w2.WriteAll(data2)
	w2.Flush()
	out = append(out, csv.NewReader(buf2))

	return out
}

func getMockBankStatementDataTcInvalidAmount() []*csv.Reader {
	out := make([]*csv.Reader, 0)

	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	data := [][]string{
		{"unique_identifier", "amount", "date"},
		{"ABC202501010396818", "-33160.38", "2025-01-01"},
	}
	w.WriteAll(data)
	w.Flush()
	out = append(out, csv.NewReader(buf))

	buf2 := &bytes.Buffer{}
	w2 := csv.NewWriter(buf2)
	data2 := [][]string{
		{"unique_identifier", "amount", "date"},
		{"DKU202501011536751", "84844.95", "2025-01-01"},
		{"DKU202501011726759", "92854.13", "2025-01-01"},
	}
	w2.WriteAll(data2)
	w2.Flush()
	out = append(out, csv.NewReader(buf2))

	return out
}

func getMockInternalTrxDataTcNotFound() *csv.Reader {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	data := [][]string{
		{"trxId", "amount", "type", "transactionTime"},
		{"ABC202501010396818", "33160.38", "DEBIT", "2025-01-01 00:00:06"},
		{"DKU202501011726759", "92854.13", "CREDIT", "2025-01-01 00:00:33"},
	}
	w.WriteAll(data)
	w.Flush()

	return csv.NewReader(buf)
}

func getMockBankStatementDataTcNotFound() []*csv.Reader {
	out := make([]*csv.Reader, 0)

	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	data := [][]string{
		{"unique_identifier", "amount", "date"},
		{"ABC202501010396818", "-33160.38", "2025-01-01"},
	}
	w.WriteAll(data)
	w.Flush()
	out = append(out, csv.NewReader(buf))

	buf2 := &bytes.Buffer{}
	w2 := csv.NewWriter(buf2)
	data2 := [][]string{
		{"unique_identifier", "amount", "date"},
		{"DKU202501011536751", "84844.94", "2025-01-01"},
	}
	w2.WriteAll(data2)
	w2.Flush()
	out = append(out, csv.NewReader(buf2))

	return out
}
