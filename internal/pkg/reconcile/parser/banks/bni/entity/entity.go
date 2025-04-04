package entity

import (
	"math"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

type CSVBankTrxData struct {
	BNIUniqueIdentifier string  `csv:"BNIUniqueIdentifier"`
	BNIDate             string  `csv:"BNIDate"`
	BNIBank             string  `csv:"-"`
	BNIAmount           float64 `csv:"BNIAmount"`
}

func (u *CSVBankTrxData) GetUniqueIdentifier() string {
	return u.BNIUniqueIdentifier
}

func (u *CSVBankTrxData) GetDate() string {
	return u.BNIDate
}

func (u *CSVBankTrxData) GetAmount() float64 {
	return u.BNIAmount
}

func (u *CSVBankTrxData) GetAbsAmount() float64 {
	return math.Abs(u.BNIAmount)
}

func (u *CSVBankTrxData) GetType() banks.TrxType {
	if u.BNIAmount <= 0 {
		return banks.DEBIT
	}

	return banks.CREDIT
}

func (u *CSVBankTrxData) GetBank() string {
	return u.BNIBank
}

func (u *CSVBankTrxData) ToBankTrxData() (returnData *banks.BankTrxData, err error) {
	t, e := time.Parse("2006-01-02", u.BNIDate)
	if e != nil {
		return nil, e
	}

	return &banks.BankTrxData{
		UniqueIdentifier: u.BNIUniqueIdentifier,
		Date:             t,
		Type:             u.GetType(),
		Bank:             u.BNIBank,
		FilePath:         "",
		Amount:           u.GetAbsAmount(),
	}, nil
}
