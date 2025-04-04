package cfs

import (
	"github.com/spf13/afero"
)

type Fs struct {
	LocalStorageFs afero.Fs
}

func NewFs(localStorageFs afero.Fs) (rd *Fs) {
	rd = &Fs{
		localStorageFs,
	}

	return
}
