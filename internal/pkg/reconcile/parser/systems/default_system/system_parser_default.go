package default_system

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/oprekable/bank-reconcile/internal/pkg/reconcile/parser/systems"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/log"
)

type CSVSystemTrxData struct {
	TrxID           string  `csv:"TrxID"`
	TransactionTime string  `csv:"TransactionTime"`
	Type            string  `csv:"Type"`
	Amount          float64 `csv:"Amount"`
}

func (u *CSVSystemTrxData) GetTrxID() string {
	return u.TrxID
}

func (u *CSVSystemTrxData) GetTransactionTime() string {
	return u.TransactionTime
}

func (u *CSVSystemTrxData) GetAmount() float64 {
	return math.Abs(u.Amount)
}

func (u *CSVSystemTrxData) GetType() systems.TrxType {
	return systems.TrxType(u.Type)
}
func (u *CSVSystemTrxData) ToSystemTrxData() (returnData *systems.SystemTrxData, err error) {
	t, e := time.Parse("2006-01-02 15:04:05", u.TransactionTime)
	if e != nil {
		return nil, e
	}

	return &systems.SystemTrxData{
		TrxID:           u.TrxID,
		TransactionTime: t,
		Type:            systems.TrxType(u.Type),
		FilePath:        "",
		Amount:          u.Amount,
	}, nil
}

type SystemParser struct {
	dataStruct        systems.SystemTrxDataInterface
	csvReader         *csv.Reader
	poolSystemTrxData *sync.Pool
	parser            systems.SystemParserType
	isHaveHeader      bool
}

var _ systems.ReconcileSystemData = (*SystemParser)(nil)

func NewSystemParser(
	dataStruct systems.SystemTrxDataInterface,
	csvReader *csv.Reader,
	isHaveHeader bool,
) (*SystemParser, error) {
	if csvReader == nil || dataStruct == nil {
		return nil, errors.New("csvReader or dataStruct is nil")
	}

	return &SystemParser{
		dataStruct:   dataStruct,
		parser:       systems.DefaultSystemParser,
		csvReader:    csvReader,
		isHaveHeader: isHaveHeader,
		poolSystemTrxData: &sync.Pool{
			New: func() interface{} {
				return &systems.SystemTrxData{}
			},
		},
	}, nil
}

func (d *SystemParser) ToSystemTrxData(ctx context.Context, filePath string) (returnData []*systems.SystemTrxData, err error) {
	var dec *csvutil.Decoder
	defer func() {
		if r := recover(); r != nil {
			errRecovery := fmt.Errorf("recovered from panic: %s", r)
			log.AddErr(ctx, errRecovery)
			return
		}
	}()

	if d.isHaveHeader {
		dec, err = csvutil.NewDecoder(d.csvReader)
		if err != nil || dec == nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	} else {
		header, _ := csvutil.Header(d.dataStruct, "csv")
		dec, err = csvutil.NewDecoder(d.csvReader, header...)
		if err != nil {
			log.AddErr(ctx, err)
			return nil, err
		}
	}

	for {
		originalData := d.dataStruct
		err = dec.Decode(originalData)
		if err != nil {
			break
		}

		//nolint:all
		//lint:ignore SA4006 sync.pool pattern just like this
		ptrSystemTrxData := d.poolSystemTrxData.Get().(*systems.SystemTrxData)
		ptrSystemTrxData, err = originalData.ToSystemTrxData()
		d.poolSystemTrxData.Put(ptrSystemTrxData)

		if err != nil {
			log.AddErr(ctx, err)
			continue
		}

		ptrSystemTrxData.FilePath = filePath
		returnData = append(returnData, ptrSystemTrxData)
	}

	if err == io.EOF {
		err = nil
	}

	return
}
