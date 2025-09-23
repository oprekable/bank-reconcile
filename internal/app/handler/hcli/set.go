package hcli

import "github.com/google/wire"

func ProviderHandlers() []Handler {
	return handlers
}

var Set = wire.NewSet(
	ProviderHandlers,
)
