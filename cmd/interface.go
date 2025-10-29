package cmd

import "github.com/spf13/cobra"

type MetaData struct {
	Usage string
	Short string
	Long  string
}

//go:generate mockery --name "Cmd" --output "./_mock" --outpkg "_mock"
type Cmd interface {
	Init(metaData *MetaData) *cobra.Command
	Runner(cCmd *cobra.Command, args []string) (er error)
	PersistentPreRunner(cCmd *cobra.Command, args []string) (er error)
	Example() string
}
