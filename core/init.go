package core

import (
	"github.com/urfave/cli/v2"
	"netspy/core/log"
	"netspy/core/misc"
)

func Init(c *cli.Context) {
	log.InitLog(c)
	misc.RecEnvInfo()
}
