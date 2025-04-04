package entity

import (
	"math"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

type CSVBankTrxData struct {
	DefaultUniqueIdentifier string  `csv:"UniqueIdentifier"`
	DefaultDate             string  `csv:"Date"`
	DefaultBank             string  `csv:"-"`
	DefaultAmount           float64 `csv:"Amount"`
}

func (u *CSVBankTrxData) GetUniqueIdentifier() string {
	return u.DefaultUniqueIdentifier
}

func (u *CSVBankTrxData) GetDate() string {
	return u.DefaultDate
}

func (u *CSVBankTrxData) GetAmount() float64 {
	return u.DefaultAmount
}

func (u *CSVBankTrxData) GetAbsAmount() float64 {
	return math.Abs(u.DefaultAmount)
}

func (u *CSVBankTrxData) GetType() banks.TrxType {
	if u.DefaultAmount <= 0 {
		return banks.DEBIT
	}

	return banks.CREDIT
}

func (u *CSVBankTrxData) GetBank() string {
	return u.DefaultBank
}

func (u *CSVBankTrxData) ToBankTrxData() (returnData *banks.BankTrxData, err error) {
	t, e := time.Parse("2006-01-02", u.DefaultDate)
	if e != nil {
		return nil, e
	}

	return &banks.BankTrxData{
		UniqueIdentifier: u.DefaultUniqueIdentifier,
		Date:             t,
		Type:             u.GetType(),
		Bank:             u.DefaultBank,
		FilePath:         "",
		Amount:           u.GetAbsAmount(),
	}, nil
}
