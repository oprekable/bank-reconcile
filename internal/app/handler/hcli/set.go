package hcli

import "github.com/google/wire"

func ProviderHandlers() []Handler {
	return Handlers
}

var Set = wire.NewSet(
	ProviderHandlers,
)
