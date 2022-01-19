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
)

var (
	thread int
	path   string
	force  bool
	mutex  = new(sync.Mutex)
)

func goSpy(ips [][]string, check func(ip string) bool) []string {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	var online []string
	var wg sync.WaitGroup
	var ipc = make(chan []string)
	go func() {
		for _, ipg := range ips {
			ipc <- ipg
		}
		defer close(ipc)
	}()

	for i := 0; i < thread; i++ {
		wg.Add(1)
		go func(ipc chan []string) {
			defer wg.Done()
			for ipg := range ipc {
				for _, ip := range ipg {
					if check(ip) {
						online = append(online, ip)
						Log.Debugf("%s alive", ip)
						Log.Printf("%s/24", ip)
						mutex.Lock()
						_, err := file.WriteString(fmt.Sprintf("%s/24\n", ip))
						if err != nil {
							Log.Error(err.Error())
						}
						mutex.Unlock()
						// 发现段内一个IP存活表示该段存活 不再检查该段
						if !force {
							println("now")
							break
						}
					} else {
						Log.Debugf("%s dead", ip)
						continue
					}
				}
			}
		}(ipc)
	}
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

func getAllCIDR(cidrs, keywords []string) []string {
	var all []string
	if cidrs == nil {
		for _, keyword := range keywords {
			if keyword == "192" {
				all = append(all, "192.168.0.0/16")
			}
			if keyword == "172" {
				all = append(all, "172.16.0.0/12")
			}
			if keyword == "10" {
				all = append(all, "10.0.0.0/8")
			}
		}
		return all
	}
	for _, cidr := range cidrs {
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			Log.Error(err)
			continue
		}
		all = append(all, cidr)
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
	keywords := c.StringSlice("net")
	cidrs := c.StringSlice("cidr")
	allcidr := getAllCIDR(cidrs, keywords)
	Log.Debugf("all cidr %v", allcidr)
	netips := getNetIPS(allcidr)
	count := c.Int("random")
	ips := GenIPS(netips, number, count)
	goSpy(ips, check)
}
