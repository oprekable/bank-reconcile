package process

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/oprekable/bank-reconcile/internal/pkg/utils/filepathhelper"
)

var Usage = "process"
var Aliases = []string{"pr", "p"}
var Short = "Process reconciliation data"
var Long = "Process reconciliation data of System Transactions and Bank Transactions"

var workDir = filepathhelper.GetWorkDir(filepathhelper.SystemCalls{})
var nowDateString = time.Now().Format("2006-01-02")

var Example = fmt.Sprintf(
	"%s --systemtrxpath=%s --banktrxpath=%s --reportpath==%s --listbank=bca,bni,mandiri,bri,danamon --from=%s --to=%s",
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
	filepath.Join(
		workDir,
		"report",
	),
	nowDateString,
	nowDateString,
)
