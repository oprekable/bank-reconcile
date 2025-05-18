package cconfig

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/creasty/defaults"
	godotenvFS "github.com/driftprogramming/godotenv"
	"github.com/joho/godotenv"
	"github.com/oprekable/bank-reconcile/internal/app/config"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	// TZ ...
	TZ string = "TZ"
)

type ConfigPaths []string
type WorkDirPath string
type AppName string
type TimeZone string
type TimeOffset int

type Config struct {
	*config.Data
	timeLocation *time.Location
	appName      AppName
	workDirPath  WorkDirPath
	timeZone     TimeZone
	timeOffset   TimeOffset
}

func initTimeZone(tzArgs TimeZone) (tz string, loc *time.Location, offset int, err error) {
	tz = os.Getenv(TZ)
	if tz == "" {
		err = os.Setenv(TZ, string(tzArgs))
		if err != nil {
			return
		}

		tz = string(tzArgs)
	}

	tzString, offset1 := time.Now().Zone()
	loc, err = time.LoadLocation(os.Getenv(TZ))
	if err != nil {
		return tzString, time.Local, offset1, nil
	}

	_, offset2 := time.Now().In(loc).Zone()

	offset = offset1

	if offset1 != offset2 {
		time.Local = loc
		offset = offset2
	}

	return tz, loc, offset, err
}

func initWorkDirPath() WorkDirPath {
	w, _ := os.UserHomeDir()
	if ex, er := os.Executable(); er == nil {
		w = filepath.Dir(ex)
	}

	return WorkDirPath(w)
}

func fromFS(embedFS fs.ReadFileFS, patterns []string, conf interface{}) (err error) {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("toml")

	var matches []string
	for i := range patterns {
		if matches, err = fs.Glob(embedFS, patterns[i]); err == nil {
			var fileData []byte
			for i2 := range matches {
				if fileData, err = embedFS.ReadFile(matches[i2]); err == nil {
					_ = viper.MergeConfig(bytes.NewReader(fileData))
				}
			}
		}
	}

	return viper.Unmarshal(conf)
}

func fromFiles(patterns []string, conf interface{}, afs afero.Fs) (err error) {
	viper.SetFs(afs)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("toml")

	var matches []string
	for i := range patterns {
		if matches, err = afero.Glob(afs, patterns[i]); err == nil {
			for i2 := range matches {
				if _, err = afs.Stat(matches[i2]); err == nil {
					viper.SetConfigFile(matches[i2])
					_ = viper.MergeInConfig()
				}
			}
		}
	}

	return viper.Unmarshal(conf)
}

func NewConfig(ctx context.Context, embedFS *embed.FS, afs afero.Fs, configPaths ConfigPaths, appName AppName, tzArgs TimeZone) (rd *Config, err error) {
	rd = &Config{}

	type InitTZStruct struct {
		Loc    *time.Location
		Tz     string
		Offset int
	}

	var cfg config.Data
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			d := InitTZStruct{}
			d.Tz, d.Loc, d.Offset, e = initTimeZone(tzArgs)
			return d, e
		},
		// Set env from embedFS file
		func(c context.Context, i interface{}) (r interface{}, e error) {
			d := i.(InitTZStruct)
			rd.workDirPath = initWorkDirPath()
			rd.timeZone = TimeZone(d.Tz)
			rd.timeLocation = d.Loc
			rd.timeOffset = TimeOffset(d.Offset)

			fileEnvPath := "embeds/envs/.env"
			return nil, godotenvFS.Load(*embedFS, fileEnvPath)
		},
		// Set env from regular file
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			if _, er := afs.Stat(string(rd.workDirPath) + "/params/.env"); er == nil {
				_ = godotenv.Overload(string(rd.workDirPath) + "/params/.env")
			}

			return
		},
		// Load config from embed files
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			cPFS := append(configPaths, "embeds/params/*.toml")
			return nil, fromFS(embedFS, cPFS, &cfg)
		},
		// Load config from regular files
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			cP := append(configPaths, fmt.Sprintf("%s/params/*.toml", rd.workDirPath))
			return nil, fromFiles(cP, &cfg, afs)
		},
		func(c context.Context, _ interface{}) (r interface{}, e error) {
			return nil, defaults.Set(&cfg)
		},
	)

	rd.Data = &cfg
	rd.appName = appName

	return rd, err
}

func (c *Config) GetWorkDirPath() WorkDirPath {
	return c.workDirPath
}

func (c *Config) GetTimeLocation() *time.Location {
	return c.timeLocation
}

func (c *Config) GetTimeOffset() TimeOffset {
	return c.timeOffset
}

func (c *Config) GetTimeZone() TimeZone {
	return c.timeZone
}

func (c *Config) GetAppName() AppName {
	return c.appName
}
