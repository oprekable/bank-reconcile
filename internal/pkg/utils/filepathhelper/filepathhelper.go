package filepathhelper

import (
	"os"
	"path/filepath"
)

type ISystemCalls interface {
	Getwd() (string, error)
	Executable() (string, error)
	FilepathDir(path string) string
}

type SystemCalls struct{}

func (SystemCalls) Getwd() (string, error) {
	return os.Getwd()
}

func (SystemCalls) Executable() (string, error) {
	return os.Executable()
}

func (SystemCalls) FilepathDir(path string) string {
	return filepath.Dir(path)
}

func GetWorkDir(sc ISystemCalls) string {
	wDirPath, err := sc.Getwd()
	if err != nil {
		ex, er := sc.Executable()
		if er == nil {
			wDirPath = sc.FilepathDir(ex)
		}
	}

	return wDirPath
}
