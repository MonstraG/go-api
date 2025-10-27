package version

import (
	"go-api/infrastructure/myLog"
	"runtime/debug"
)

var AppVersion string

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		myLog.Fatal.Logf("Failed to read build info")
		return
	}

	AppVersion = buildInfo.Main.Version

	myLog.Info.Logf("Version: %s", AppVersion)
}
