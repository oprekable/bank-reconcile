package version

import (
	"fmt"
	"io"
	"os"

	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/versionhelper"
	"github.com/oprekable/bank-reconcile/variable"

	"github.com/spf13/cobra"
)

var versionWriter io.Writer = os.Stdout

func Runner(_ *cobra.Command, _ []string) (er error) {
	start(versionWriter)
	return nil
}

func start(w io.Writer) {
	defer func() {
		atexit.AtExit()
	}()

	atexit.Add(
		func() {
			shutdown(w)
		},
	)

	_, _ = fmt.Fprintln(w, "App\t\t\t:", variable.AppName)
	_, _ = fmt.Fprintln(w, "Desc\t\t:", variable.AppDescLong)
	_, _ = fmt.Fprintln(w, "Build Date\t:", variable.BuildDate)
	_, _ = fmt.Fprintln(w, "Git Commit\t:", variable.GitCommit)
	_, _ = fmt.Fprintln(w, "Version\t\t:", versionhelper.GetVersion(variable.Version))
	_, _ = fmt.Fprintln(w, "environment\t:", variable.Environment)
	_, _ = fmt.Fprintln(w, "Go Version\t:", variable.GoVersion)
	_, _ = fmt.Fprintln(w, "OS / Arch\t:", variable.OsArch)
}

func shutdown(w io.Writer) {
	_, _ = fmt.Fprintln(w, "\n-#-")
}
