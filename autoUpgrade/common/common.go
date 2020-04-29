package common

import (
	"runtime"
)

//SysType is the value of windows or linux or others
const SysType = runtime.GOOS

//UpgradeLog is log folder of autoUpgrade
const UpgradeLog = "upgradeLog"

