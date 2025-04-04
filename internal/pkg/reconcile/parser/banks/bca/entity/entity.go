package entity

import (
	"math"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
)

type CSVBankTrxData struct {
	BCAUniqueIdentifier string  `csv:"BCAUniqueIdentifier"`
	BCADate             string  `csv:"BCADate"`
	BCABank             string  `csv:"-"`
	BCAAmount           float64 `csv:"BCAAmount"`
}

func (u *CSVBankTrxData) GetUniqueIdentifier() string {
	return u.BCAUniqueIdentifier
}

func (u *CSVBankTrxData) GetDate() string {
	return u.BCADate
}

func (u *CSVBankTrxData) GetAmount() float64 {
	return u.BCAAmount
}

func (u *CSVBankTrxData) GetAbsAmount() float64 {
	return math.Abs(u.BCAAmount)
}

func (u *CSVBankTrxData) GetType() banks.TrxType {
	if u.BCAAmount <= 0 {
		return banks.DEBIT
	}

	return banks.CREDIT
}

func (u *CSVBankTrxData) GetBank() string {
	return u.BCABank
}

func (u *CSVBankTrxData) ToBankTrxData() (returnData *banks.BankTrxData, err error) {
	t, e := time.Parse("2006-01-02", u.BCADate)
	if e != nil {
		return nil, e
	}

	return &banks.BankTrxData{
		UniqueIdentifier: u.BCAUniqueIdentifier,
		Date:             t,
		Type:             u.GetType(),
		Bank:             u.BCABank,
		FilePath:         "",
		Amount:           u.GetAbsAmount(),
	}, nil
}
