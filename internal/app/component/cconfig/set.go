package cconfig

import (
	"github.com/google/wire"
	"github.com/spf13/afero"
)

type ConfigFS afero.Fs

var Set = wire.NewSet(
	wire.InterfaceValue(new(ConfigFS), afero.NewOsFs()),
	NewConfig,
)
