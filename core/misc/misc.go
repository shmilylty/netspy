package misc

import (
	"net"
	. "netspy/core/log"
	"runtime"
)

func RecEnvInfo() {
	Log.Debugf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// CalcBcstIP 计算广播地址(网段最后一个IP)
func CalcBcstIP(c *net.IPNet) net.IP {
	mask := c.Mask
	bcst := make(net.IP, len(c.IP))
	copy(bcst, c.IP)
	for i := 0; i < len(mask); i++ {
		ipIdx := len(bcst) - i - 1
		bcst[ipIdx] = c.IP[ipIdx] | ^mask[len(mask)-i-1]
	}
	return bcst
}
