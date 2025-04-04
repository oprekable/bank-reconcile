package cfs

import (
	"github.com/google/wire"
	"github.com/spf13/afero"
)

func ProviderCFs(fsType afero.Fs) *Fs {
	return NewFs(
		fsType,
	)
}

var Set = wire.NewSet(
	ProviderCFs,
)
