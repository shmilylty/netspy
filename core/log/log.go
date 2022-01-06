package log

import (
	"github.com/kataras/golog"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

var Log = golog.New()

func InitLog(c *cli.Context) {
	if c.Bool("debug") {
		Log.SetLevel("debug")
	}
	if c.Bool("silent") {
		Log.SetLevel("error")
		Log.SetTimeFormat("")
	}
	logFile, err := os.OpenFile("netspy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}
