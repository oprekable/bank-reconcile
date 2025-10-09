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
	buildInfo, _ := debug.ReadBuildInfo()
	return getVersionLogic(buildInfo, version, buildDate, commitHash, environment)
}

func getVersionLogic(buildInfo *debug.BuildInfo, version, buildDate, commitHash, environment string) (returnData VersionStruct) {
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

	if returnData.Version == Snapshot {
		returnData.Version = buildInfo.Main.Version
	}

	for i := range buildInfo.Settings {
		switch buildInfo.Settings[i].Key {
		case "vcs.revision":
			{
				returnData.CommitHash = buildInfo.Settings[i].Value
			}
		case "vcs.time":
			{
				returnData.BuildDate = buildInfo.Settings[i].Value
			}
		}
	}

	return returnData
}
