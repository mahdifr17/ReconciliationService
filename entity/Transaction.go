package entity

type Transaction struct {
	TrxId string
	Amount float64
	Type TransactionType
	TransactionTime time.Time // datetime format
}

func (t Transaction) ReadFromCsv(input string) (err error) {
	// input must be trx_id, amount, type, transaction_time
	tmp := strings.Split(input, ",")
    t.TrxId = tmp[0]
    t.Amount, err = strconv.ParseFloat(tmp[1], 64)
	if err != nil {
		return err
	}
	t.Type = TransactionType()
	t.TransactionTime, err = time.Parse(time.DateTime, tmp[3])
	if err != nil {
		return err
	}
}