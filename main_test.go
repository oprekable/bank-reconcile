package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/oprekable/bank-reconcile/cmd/process"
	"github.com/oprekable/bank-reconcile/cmd/sample"
	"github.com/oprekable/bank-reconcile/cmd/version"
	"github.com/oprekable/bank-reconcile/variable"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestMainApp(t *testing.T) {
	origExit := exitFunc
	origArgs := os.Args
	defer func() {
		exitFunc = origExit
		os.Args = origArgs
	}()

	t.Run("Success exit code 0", func(t *testing.T) {
		var capturedCode int
		exitFunc = func(code int) {
			capturedCode = code
		}

		os.Args = []string{
			variable.AppName,
			version.Usage,
		}

		main()
		assert.Equal(t, 0, capturedCode)
	})

	t.Run("Error exit code 1", func(t *testing.T) {
		var capturedCode int
		exitFunc = func(code int) {
			capturedCode = code
		}

		os.Args = []string{
			variable.AppName,
			"--invalid-flag",
		}

		main()
		assert.Equal(t, 1, capturedCode)
	})
}

func TestMainLogic(t *testing.T) {
	var outPutWriter io.Writer = new(bytes.Buffer)
	var errWriter io.Writer = new(bytes.Buffer)

	t.Run("Running app `version` command returns 0", func(t *testing.T) {
		os.Args = []string{
			variable.AppName,
			version.Usage,
		}

		exitCode := run(outPutWriter, errWriter)
		assert.Equal(t, 0, exitCode)
	})

	t.Run("Running app `sample` command returns 0", func(t *testing.T) {
		os.Args = []string{
			variable.AppName,
			sample.Usage,
		}

		exitCode := run(outPutWriter, errWriter)
		assert.Equal(t, 0, exitCode)
	})

	t.Run("Running app `process` command returns 0", func(t *testing.T) {
		os.Args = []string{
			variable.AppName,
			process.Usage,
		}

		exitCode := run(outPutWriter, errWriter)
		assert.Equal(t, 0, exitCode)
	})

	t.Run("Running app with invalid command returns 1 (Execute error branch)", func(t *testing.T) {
		os.Args = []string{
			variable.AppName,
			"invalid-subcommand",
		}

		exitCode := run(outPutWriter, errWriter)
		assert.Equal(t, 1, exitCode)
	})

	t.Run("Running app with unknown flag returns 1 (Execute error branch)", func(t *testing.T) {
		os.Args = []string{
			variable.AppName,
			"--unknown-flag",
		}

		exitCode := run(outPutWriter, errWriter)
		assert.Equal(t, 1, exitCode)
	})
}
