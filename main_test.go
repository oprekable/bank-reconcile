package main

import (
	"github.com/oprekable/bank-reconcile/cmd/process"
	"github.com/oprekable/bank-reconcile/cmd/sample"
	"github.com/oprekable/bank-reconcile/cmd/version"
	"github.com/oprekable/bank-reconcile/variable"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestMainApp(t *testing.T) {
	t.Log("Running app `version` command")
	os.Args = []string{
		variable.AppName,
		version.Usage,
	}

	main()

	t.Log("Running app `sample` command")
	os.Args = []string{
		variable.AppName,
		sample.Usage,
	}

	main()

	t.Log("Running app `process` command")
	os.Args = []string{
		variable.AppName,
		process.Usage,
	}

	main()
}
