package spy

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"netspy/core/ips"
	. "netspy/core/log"
	"os"
	"sync"
)

var (
	thread int
	path   string
	m      = new(sync.Mutex)
)

func goSpy(ips []string, check func(ip string) bool) []string {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	var online []string
	var wg sync.WaitGroup
	var ipc = make(chan string)
	go func() {
		for _, ip := range ips {
			ipc <- ip
		}
		defer close(ipc)
	}()

	for i := 0; i < thread; i++ {
		wg.Add(1)
		go func(ipc chan string) {
			defer wg.Done()
			for ip := range ipc {
				if check(ip) {
					online = append(online, ip)
					Log.Infof("%s/24", ip)
					m.Lock()
					_, err := file.WriteString(fmt.Sprintf("%s/24\n", ip))
					if err != nil {
						Log.Error(err.Error())
					}
					m.Unlock()
				} else {
					continue
				}
			}
		}(ipc)
	}
	wg.Wait()
	return online
}

func Spy(c *cli.Context, check func(ip string) bool) {
	path = c.Path("output")
	thread = c.Int("thread")
	keyword := c.String("net")
	number := c.StringSlice("number")
	var gateways, all []string
	if keyword == "all" || keyword == "192" {
		Log.Info("start to spy 192.168.0.0/16")
		gateways = ips.GenIps("192.168.0.1", number, "b")
		all = append(all, gateways...)
		goSpy(gateways, check)
	}

	if keyword == "all" || keyword == "172" {
		Log.Info("start to spy 172.16.0.0/12")
		gateways = ips.GenIps("172.16.0.0", number, "172")
		all = append(all, gateways...)
		goSpy(gateways, check)
	}

	if keyword == "all" || keyword == "10" {
		Log.Info("start to spy 10.0.0.0/8")
		gateways = ips.GenIps("10.0.0.1", number, "a")
		all = append(all, gateways...)
		goSpy(gateways, check)
	}
}
