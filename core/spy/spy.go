package spy

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"net"
	. "netspy/core/log"
	"netspy/core/misc"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	thread int
	path   string
	force  bool
)

func goSpy(ips [][]string, check func(ip string) bool) []string {
	var online []string
	var wg sync.WaitGroup
	var ipc = make(chan []string, 10000)
	var onc = make(chan string, 1000)
	var count int32

	if ips == nil {
		return online
	}
	go func() {
		for _, ipg := range ips {
			ipc <- ipg
		}
		defer close(ipc)
	}()

	// 探测协程
	for i := 0; i < thread; i++ {
		wg.Add(1)
		go func(ipc chan []string) {
			for ipg := range ipc {
				for _, ip := range ipg {
					if check(ip) {
						online = append(online, ip)
						Log.Debugf("%s alive", ip)
						Log.Printf("%s/24", ip)
						s := fmt.Sprintf("%s/24\n", ip)
						onc <- s
						// 发现段内一个IP存活表示该段存活 不再检查该段
						if !force {
							break
						}
					} else {
						Log.Debugf("%s dead", ip)
						continue
					}
				}
				atomic.AddInt32(&count, int32(len(ipg)))
			}
			defer wg.Done()
		}(ipc)
	}

	// 保存协程
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Error(err.Error())
	}
	defer file.Close()

	go func(onc chan string) {
		for s := range onc {
			_, err := file.WriteString(s)
			if err != nil {
				Log.Error(err.Error())
			}
		}
	}(onc)

	// 统计协程
	num := len(ips[0])
	wg.Add(1)
	go func() {
		all := float64(len(ips) * num)
		i := 0
		for {
			time.Sleep(10 * time.Second)
			i += 1
			spied := float64(count)
			speed := float64(count) / (float64(i) * 10)
			remain := (all - spied) / speed
			Log.Infof("all: %d spied: %d ratio: %.2f speed: %.2f it/s remain: %.0fs",
				all, spied, spied/all, speed, remain)
			if all == spied {
				wg.Done()
				break
			}
		}
	}()

	wg.Wait()
	return online
}

func setThread(i int) int {
	if i == 0 {
		return runtime.NumCPU() * 20
	}
	return i
}

func genNetIP(start, end net.IP) []net.IP {
	var netip []net.IP
	// 10.0.0.0 - 10.0.0.255 情况
	// 10.0.0.0 - 10.0.10.255 情况
	if start[0] == end[0] && start[1] == end[1] {
		for k := start[2]; k <= end[2]; k++ {
			// 放入循环是为了每次循环创建内存地址不同的新IP
			ip := make(net.IP, len(start))
			// 深拷贝
			copy(ip, start)
			ip[2] = k
			netip = append(netip, ip)
			if k == 255 {
				break
			}
		}
	}
	// 10.0.0.0 - 10.10.255.255 情况
	if start[0] == end[0] && start[1] != end[1] {
		for j := start[1]; j <= end[1]; j++ {
			for k := start[2]; k <= end[2]; k++ {
				ip := make(net.IP, len(start))
				copy(ip, start)
				ip[1] = j
				ip[2] = k
				netip = append(netip, ip)
				if k == 255 {
					break
				}
			}
			if j == 255 {
				break
			}
		}
	}

	// 10.0.0.0 - 20.255.255.255 这种情况不一定存在
	if start[0] != end[0] {
		for i := start[0]; i <= end[0]; i++ {
			for j := start[1]; j <= end[1]; j++ {
				for k := start[2]; k <= end[2]; k++ {
					ip := make(net.IP, len(start))
					copy(ip, start)
					ip[0] = i
					ip[1] = j
					ip[2] = k
					netip = append(netip, ip)
					if k == 255 {
						break
					}
				}
				if j == 255 {
					break
				}
			}
			if i == 255 {
				break
			}
		}
	}
	return netip
}

func getNetIPS(cidrs []string) []net.IP {
	var netips []net.IP
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			Log.Fatal(err)
		}
		start := ipnet.IP
		end := misc.CalcBcstIP(ipnet)
		Log.Infof("%v is from %v to %v", cidr, start, end)
		netip := genNetIP(start, end)
		netips = append(netips, netip...)
	}
	return netips
}

func genAllCIDR() []string {
	var all []string
	c := [9]int{1, 32, 64, 96, 128, 160, 192, 224, 255}
	for i := 1; i <= 255; i++ {
		for j := 1; j <= 255; j++ {
			for _, k := range c {
				cidr := fmt.Sprintf("%v.%v.%v.0/24", i, j, k)
				all = append(all, cidr)
			}
		}
	}
	return all
}

func mergeCIDR(cidrs []string, special bool) []string {
	var all []string
	for _, cidr := range cidrs {
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			Log.Error(err)
			continue
		}
		all = append(all, cidr)
	}
	if all != nil {
		return all
	}
	if all == nil {
		c := []string{"192.168.0.0/16", "172.16.0.0/12", "10.0.0.0/8"}
		all = append(all, c...)
	}
	if special {
		if misc.IsPureIntranet() {
			Log.Debug("the current network environment is pure intranet")
			all = genAllCIDR()
		} else {
			c := []string{"100.64.0.0/10", "59.192.0.0/10", "3.1.0.0/10"}
			all = append(all, c...)
		}
	}
	return all
}

func checkEndNum(nums []string) []int {
	var tail []int
	for _, s := range nums {
		i, err := strconv.Atoi(s)
		if err != nil {
			Log.Error(err)
			continue
		}
		if i >= 0 && i <= 255 {
			tail = append(tail, i)
		}
	}
	return tail
}

func Spy(c *cli.Context, check func(ip string) bool) {
	thread = setThread(c.Int("thread"))
	Log.Debugf("%v threads", thread)
	path = c.Path("output")
	Log.Debugf("save path: %v", path)
	force = c.Bool("force")
	number := checkEndNum(c.StringSlice("end"))
	special := c.Bool("special")
	cidrs := c.StringSlice("cidr")
	allcidr := mergeCIDR(cidrs, special)
	Log.Debugf("all cidr %v", allcidr)
	netips := getNetIPS(allcidr)
	count := c.Int("random")
	ips := GenIPS(netips, number, count)
	goSpy(ips, check)
}
