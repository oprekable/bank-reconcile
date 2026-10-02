package entity

import (
	"math"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

type CSVBankTrxData struct {
	BSIUniqueIdentifier string  `csv:"BSIUniqueIdentifier"`
	BSIDate             string  `csv:"BSIDate"`
	BSIBank             string  `csv:"-"`
	BSIAmount           float64 `csv:"BSIAmount"`
}

func (u *CSVBankTrxData) GetUniqueIdentifier() string {
	return u.BSIUniqueIdentifier
}

func (u *CSVBankTrxData) GetDate() string {
	return u.BSIDate
}

func (u *CSVBankTrxData) GetAmount() float64 {
	return u.BSIAmount
}

func (u *CSVBankTrxData) GetAbsAmount() float64 {
	return math.Abs(u.BSIAmount)
}

func (u *CSVBankTrxData) GetType() banks.TrxType {
	if u.BSIAmount <= 0 {
		return banks.DEBIT
	}

	return banks.CREDIT
}

func (u *CSVBankTrxData) GetBank() string {
	return u.BSIBank
}

func (u *CSVBankTrxData) ToBankTrxData() (returnData *banks.BankTrxData, err error) {
	t, e := time.Parse("2006-01-02", u.BSIDate)
	if e != nil {
		return nil, e
	}

	return &banks.BankTrxData{
		UniqueIdentifier: u.BSIUniqueIdentifier,
		Date:             t,
		Type:             u.GetType(),
		Bank:             u.BSIBank,
		FilePath:         "",
		Amount:           u.GetAbsAmount(),
	}, nil
}
