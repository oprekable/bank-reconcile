package versionhelper

import (
	"runtime/debug"
)

type VersionStruct struct {
	Version     string
	BuildDate   string
	CommitHash  string
	Environment string
}

const (
	Snapshot = "SNAPSHOT"
	Default  = "default"
)

func GetVersion(version, buildDate, commitHash, environment string) (returnData VersionStruct) {
	returnData.Version = version
	returnData.BuildDate = buildDate
	returnData.CommitHash = commitHash
	returnData.Environment = environment

	if returnData.Version == "" {
		returnData.Version = Snapshot
	}

	if returnData.Environment == "" {
		returnData.Environment = Default
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return returnData
	}

	if returnData.Version == Snapshot {
		returnData.Version = info.Main.Version
	}

	for i := range info.Settings {
		switch info.Settings[i].Key {
		case "vcs.revision":
			{
				returnData.CommitHash = info.Settings[i].Value
			}
		case "vcs.time":
			{
				returnData.BuildDate = info.Settings[i].Value
			}
		}
	}

	return returnData
}
