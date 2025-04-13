package sample

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
)

var Usage = "sample"
var Aliases = []string{"sa", "s"}
var Short = "Generate sample reconciliation data"
var Long = "Generate sample reconciliation data of System Transactions and Bank Transactions"

var workDir = filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})
var nowDateString = time.Now().Format("2006-01-02")

var Example = fmt.Sprintf(
	"%s --systemtrxpath=%s --banktrxpath=%s --listbank=bca,bni,mandiri,bri,danamon --percentagematch=100 --amountdata=10000 --from=%s --to=%s",
	Usage,
	filepath.Join(
		workDir,
		"sample",
		"system",
	),
	filepath.Join(
		workDir,
		"sample",
		"bank",
	),
	nowDateString,
	nowDateString,
)
