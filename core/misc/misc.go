package misc

import (
	"bytes"
	"net"
	. "netspy/core/log"
	"os/exec"
	"runtime"
	"strings"
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

func IsPing(ip, times, timeout string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c",
			"ping -n "+times+" -w "+timeout+" "+ip+" && echo true || echo false")
		break
	case "linux":
		cmd = exec.Command("/bin/sh", "-c",
			"ping -c "+times+" -w "+timeout+" "+ip+" >/dev/null && echo true || echo false")
		break
	case "darwin":
		cmd = exec.Command("/bin/sh", "-c",
			"ping -c "+times+" -w "+timeout+" "+ip+" >/dev/null && echo true || echo false")
		break
	default:
		cmd = nil
	}

	var output = bytes.Buffer{}
	if cmd != nil {
		cmd.Stdout = &output
		var err = cmd.Start()
		if err != nil {
			return false
		}
		if err = cmd.Wait(); err != nil {
			return false
		} else {
			if strings.Contains(output.String(), "true") {
				return true
			} else {
				return false
			}
		}
	} else {
		return false
	}
}

func IsPureIntranet() bool {
	if IsPing("114.114.114.114", string('3'), string('3')) {
		return false
	}
	if IsPing("8.8.8.8", string('3'), string('3')) {
		return false
	}
	return true
}
