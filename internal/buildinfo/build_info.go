package buildinfo

import "runtime/debug"

var BuildInfo *debug.BuildInfo

func Version() string {
	return BuildInfo.Main.Version
}

func init() {
	var ok bool
	BuildInfo, ok = debug.ReadBuildInfo()
	if !ok {
		panic("cannot read build info")
	}
}
