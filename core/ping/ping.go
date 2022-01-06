package ping

import (
	"bytes"
	"github.com/urfave/cli/v2"
	. "netspy/core/log"
	"netspy/core/spy"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var (
	times   string
	timeout string
	goos    = runtime.GOOS
)

func check(ip string) bool {
	var cmd *exec.Cmd
	switch goos {
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

func Spy(c *cli.Context) {
	Log.Info("use ping command to spy")
	times = strconv.Itoa(c.Int("times"))
	timeout = strconv.Itoa(c.Int("timeout"))
	spy.Spy(c, check)
}
