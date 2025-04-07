package usecase

import (
	"bytes"
	"encoding/csv"
	"testing"
	"time"

	"github.com/mahdifr17/ReconciliationService/entity"
)

func TestReconciliationUsecaseImpl_ReconcileData(t *testing.T) {
	type args struct {
		csvInternalTrxData *csv.Reader
		csvBankStatements  []*csv.Reader
		startDate          time.Time
		endDate            time.Time
	}

	startDate, _ := time.Parse(time.DateOnly, "2025-01-01")
	endDate, _ := time.Parse(time.DateOnly, "2025-02-28")
	endDate = endDate.Add(86399 * time.Second) // 23:59:59

	tests := []struct {
		name string
		uc   ReconciliationUsecase
		args args
		want entity.ReconciliationResult
	}{
		{
			name: "valid trx",
			uc:   new(ReconciliationUsecaseImpl),
			args: args{getMockInternalTrxDataTcValid(), getMockBankStatementDataTcValid(), startDate, endDate},
			want: entity.ReconciliationResult{
				TotalTransactionProcessed: 3,
				TotalMatchTransaction:     3,
			},
		},
		{
			name: "invalid amount",
			uc:   new(ReconciliationUsecaseImpl),
			args: args{getMockInternalTrxDataTcValid(), getMockBankStatementDataTcInvalidAmount(), startDate, endDate},
			want: entity.ReconciliationResult{
				TotalTransactionProcessed: 3,
				TotalMatchTransaction:     2,
				TotalDiscrepancies:        0.00999999999476131,
			},
		},
		{
			name: "not found on internal",
			uc:   new(ReconciliationUsecaseImpl),
			args: args{getMockInternalTrxDataTcNotFound(), getMockBankStatementDataTcValid(), startDate, endDate},
			want: entity.ReconciliationResult{
				TotalTransactionProcessed: 3,
				TotalMatchTransaction:     2,
				ListMissingTransactionInternal: map[string][]entity.CompareResult{
					"1": {
						{
							InternalTransaction: entity.Transaction{
								TrxId: "DKU202501011536751",
							},
							Remark: errCaseNotFoundOnInternal,
						},
					},
				},
			},
		},
		{
			name: "not found on bank statement",
			uc:   new(ReconciliationUsecaseImpl),
			args: args{getMockInternalTrxDataTcValid(), getMockBankStatementDataTcNotFound(), startDate, endDate},
			want: entity.ReconciliationResult{
				TotalTransactionProcessed: 2,
				TotalMatchTransaction:     2,
				ListMissingTransactionBank: []entity.CompareResult{
					{
						InternalTransaction: entity.Transaction{
							TrxId: "DKU202501011726759",
						},
						Remark: errCaseNotFoundOnBankStatement,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.uc.ReconcileData(tt.args.csvInternalTrxData, tt.args.csvBankStatements, tt.args.startDate, tt.args.endDate)
			if got.TotalTransactionProcessed != tt.want.TotalTransactionProcessed {
				t.Errorf("ReconciliationUsecaseImpl.ReconcileData() TotalTransactionProcessed = %v, want %v", got.TotalTransactionProcessed, tt.want.TotalTransactionProcessed)
			}
			if got.TotalMatchTransaction != tt.want.TotalMatchTransaction {
				t.Errorf("ReconciliationUsecaseImpl.ReconcileData() TotalMatchTransaction = %v, want %v", got.TotalMatchTransaction, tt.want.TotalMatchTransaction)
			}
			if got.TotalDiscrepancies != tt.want.TotalDiscrepancies {
				t.Errorf("ReconciliationUsecaseImpl.ReconcileData() TotalDiscrepancies = %v, want %v", got.TotalDiscrepancies, tt.want.TotalDiscrepancies)
			}
			if len(got.ListMissingTransactionBank) != len(tt.want.ListMissingTransactionBank) {
				t.Errorf("ReconciliationUsecaseImpl.ReconcileData() ListMissingTransactionBank = %v, want %v", len(got.ListMissingTransactionBank), len(tt.want.ListMissingTransactionBank))
			}
			for k := range got.ListMissingTransactionInternal {
				if len(got.ListMissingTransactionInternal[k]) != len(tt.want.ListMissingTransactionInternal[k]) {
					t.Errorf("ReconciliationUsecaseImpl.ReconcileData() ListMissingTransactionInternal idx:%v = %v, want %v", k, len(got.ListMissingTransactionInternal[k]), len(tt.want.ListMissingTransactionInternal[k]))
				}
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
