package reconciliation

import "time"

// Reconciliation ..
type Reconciliation struct {
	FromDate                       time.Time `default:"-"    mapstructure:"from_date"`
	ToDate                         time.Time `default:"-"    mapstructure:"to_date"`
	Action                         string    `default:"-"    mapstructure:"action"`
	SystemTRXPath                  string    `default:"-"    mapstructure:"system_trx_path"`
	BankTRXPath                    string    `default:"-"    mapstructure:"bank_trx_path"`
	ReportTRXPath                  string    `default:"-"    mapstructure:"report_trx_path"`
	ListBank                       []string  `default:"-"    mapstructure:"list_bank"`
	TotalData                      int64     `default:"-"    mapstructure:"total_data"`
	PercentageMatch                int       `default:"100"  mapstructure:"percentage_match"`
	NumberWorker                   int       `default:"10"   mapstructure:"number_worker"`
	IsDeleteCurrentSampleDirectory bool      `default:"true" mapstructure:"is_delete_current_sample_directory"`
}
