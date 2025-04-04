package version

import (
	"fmt"
	"os"

	"github.com/oprekable/bank-reconcile/internal/pkg/utils/atexit"
	"github.com/oprekable/bank-reconcile/internal/pkg/utils/versionhelper"
	"github.com/oprekable/bank-reconcile/variable"

	"github.com/spf13/cobra"
)

func Runner(_ *cobra.Command, _ []string) (er error) {
	start()
	return nil
}

func start() {
	code := 0

	defer func() {
		os.Exit(code)
	}()

	defer func() {
		atexit.AtExit()
	}()

	atexit.Add(shutdown)

	fmt.Println("App\t\t:", variable.AppName)
	fmt.Println("Desc\t\t:", variable.AppDescLong)
	fmt.Println("Build Date\t:", variable.BuildDate)
	fmt.Println("Git Commit\t:", variable.GitCommit)
	fmt.Println("Version\t\t:", versionhelper.GetVersion(variable.Version))
	fmt.Println("environment\t:", variable.Environment)
	fmt.Println("Go Version\t:", variable.GoVersion)
	fmt.Println("OS / Arch\t:", variable.OsArch)
}

func shutdown() {
	fmt.Println("\n-#-")
}
