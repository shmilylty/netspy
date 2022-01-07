package misc

import (
	. "netspy/core/log"
	"runtime"
)

func RecEnvInfo() {
	Log.Debugf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
