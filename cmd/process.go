package cmd

import (
	"fmt"

	"github.com/oprekable/bank-reconcile/cmd/process"
	"github.com/oprekable/bank-reconcile/cmd/root"
	"github.com/oprekable/bank-reconcile/variable"

	"github.com/spf13/cobra"
)

var processCmd = &cobra.Command{
	Use:     process.Usage,
	Aliases: process.Aliases,
	Short:   process.Short,
	Long:    process.Long,
	Example: fmt.Sprintf(
		"%s\n",
		fmt.Sprintf("Process data \n\t%s process %s", variable.AppName, root.ProcessUsageFlags),
	),
	RunE: process.Runner,
}

func init() {
	rootCmd.AddCommand(processCmd)
}
