package sample

const (
	DEBIT = "DEBIT"
)

type Summary struct {
	TotalBankTrx   map[string]int64
	FileBankTrx    map[string]string
	FileSystemTrx  string
	TotalSystemTrx int64
}
