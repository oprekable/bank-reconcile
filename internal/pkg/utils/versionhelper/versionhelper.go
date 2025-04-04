package versionhelper

func GetVersion(version string) string {
	if version == "" {
		return "SNAPSHOT"
	}

	return version
}
