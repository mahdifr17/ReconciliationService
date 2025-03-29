package entity

type BankStatement struct {
	UniqueIdentifier string
	Amount float64 // negative value means Debit
	Date time.Time // date
}

func (bs BankStatement) ReadFromCsv(input string) (err error) {
	// input must be unique_identifier, amount, date
	tmp := strings.Split(input, ",")
    bs.UniqueIdentifier = tmp[0]
    bs.Amount, err = strconv.ParseFloat(tmp[1], 64)
	if err != nil {
		return err
	}
	bs.Date, err = time.Parse(time.DateOnly, tmp[2])
	if err != nil {
		return err
	}
}