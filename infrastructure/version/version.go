package version

import (
	"go-api/infrastructure/myLog"
	"runtime/debug"
	"time"
)

var AppVersion string
var AppBuildTime time.Time

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		myLog.Fatal.Logf("Failed to read build info")
		return
	}

	AppVersion = buildInfo.Main.Version

	for _, setting := range buildInfo.Settings {
		if setting.Key == "vcs.time" {
			appBuildTime, err := time.Parse(time.RFC3339, setting.Value)
			if err == nil {
				AppBuildTime = appBuildTime
			}
			break
		}
	}

	myLog.Info.Logf("Version: %s", AppVersion)
}
