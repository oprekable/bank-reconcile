package main

import (
	"embed"

	"github.com/oprekable/bank-reconcile/cmd"
)

//go:embed all:embeds
var embedFS embed.FS

func main() {
	cmd.Execute(&embedFS)
}
