package csvhelper

import (
	"context"
	"path/filepath"

	"github.com/aaronjan/hunch"
	"github.com/jszwec/csvutil"
	"github.com/spf13/afero"
)

func StructToCSVFile(ctx context.Context, fs afero.Fs, filePath string, structData interface{}, isDeleteDirectory bool) error {
	_, err := hunch.Waterfall(
		ctx,
		func(c context.Context, i interface{}) (interface{}, error) {
			if isDeleteDirectory {
				return nil, DeleteDirectory(c, fs, filePath)
			}

			return nil, nil
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return nil, fs.MkdirAll(filepath.Dir(filePath), 0755)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			return csvutil.Marshal(structData)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			marshal := i.([]byte)
			return nil, afero.WriteFile(fs, filePath, marshal, 0644)
		},
	)

	return err
}

func DeleteDirectory(ctx context.Context, fs afero.Fs, filePath string) (err error) {
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, i interface{}) (interface{}, error) {
			filePath = filepath.Clean(filePath)
			basePath := filepath.Dir(filePath)
			return afero.Glob(fs, basePath+"/*")
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			contents := i.([]string)
			for _, item := range contents {
				e := fs.RemoveAll(item)
				if e != nil {
					return nil, e
				}
			}
			return nil, nil
		},
	)

	return err
}
