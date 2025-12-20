package main

import (
	"io"
	"os"
	"testing"

	"github.com/oprekable/bank-reconcile/cmd/process"
	"github.com/oprekable/bank-reconcile/cmd/sample"
	"github.com/oprekable/bank-reconcile/cmd/version"
	"github.com/oprekable/bank-reconcile/variable"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestMainApp(_ *testing.T) {
	os.Args = []string{
		variable.AppName,
		version.Usage,
	}

	exitFunc = func(c int) {}

	main()
}

func TestMainLogic(t *testing.T) {
	var outPutWriter io.Writer = os.Stdout
	var errWriter io.Writer = os.Stderr

	t.Log("Running app `version` command")
	os.Args = []string{
		variable.AppName,
		version.Usage,
	}

	mainLogic(outPutWriter, errWriter)

	t.Log("Running app `sample` command")
	os.Args = []string{
		variable.AppName,
		sample.Usage,
	}

	mainLogic(outPutWriter, errWriter)

	t.Log("Running app `process` command")
	os.Args = []string{
		variable.AppName,
		process.Usage,
	}

	mainLogic(outPutWriter, errWriter)
}
