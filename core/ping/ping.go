package ping

import (
	"github.com/urfave/cli/v2"
	. "netspy/core/log"
	"netspy/core/misc"
	"netspy/core/spy"
	"strconv"
)

var (
	times   string
	timeout string
)

func check(ip string) bool {
	if misc.IsPing(ip, times, timeout) {
		return true
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
