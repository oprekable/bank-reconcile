package helper

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/banks"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"

	"github.com/jszwec/csvutil"
)

func ToBankTrxData(ctx context.Context, filePath string, isHaveHeader bool, bank string, csvReader *csv.Reader, originalData banks.BankTrxDataInterface) (returnData []*banks.BankTrxData, err error) {
	var dec *csvutil.Decoder
	defer func() {
		if r := recover(); r != nil {
			errRecovery := fmt.Errorf("recovered from panic: %s", r)
			log.AddErr(ctx, errRecovery)
			return
		}
	}()

	if isHaveHeader {
		dec, err = csvutil.NewDecoder(csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, _ := csvutil.Header(originalData, "csv")
		dec, err = csvutil.NewDecoder(csvReader, header...)
		if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	}

	for {
		err = dec.Decode(originalData)
		if err != nil {
			break
		}

		bankTrxData, er := originalData.ToBankTrxData()
		if er != nil {
			log.AddErr(ctx, er)
			continue
		}

		bankTrxData.Bank = bank
		bankTrxData.FilePath = filePath
		bankTrxData.Type = originalData.GetType()
		returnData = append(returnData, bankTrxData)
	}

	if err == io.EOF {
		err = nil
	}

	return
}
