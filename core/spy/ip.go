package spy

import (
	"fmt"
	"strings"
)

func genBClassIps(ip string, num []string) [][]string {
	var ips [][]string
	s := strings.Split(ip, ".")
	for i := 0; i < 255; i++ {
		var ipg []string
		for _, v := range num {
			// ip group
			ipg = append(ipg, fmt.Sprintf("%s.%s.%d.%s", s[0], s[1], i, v))
		}
		ips = append(ips, ipg)
	}
	return ips
}

func genAClassIps(ip string, num []string) [][]string {
	var ips [][]string
	s := strings.Split(ip, ".")
	for i := 0; i < 255; i++ {
		ip = fmt.Sprintf("%s.%d.%s.%s", s[0], i, s[2], s[3])
		ips = genBClassIps(ip, num)
	}
	return ips
}

func gen172ClassIps(ip string, num []string) [][]string {
	var ips [][]string
	s := strings.Split(ip, ".")
	for i := 16; i < 32; i++ {
		ip = fmt.Sprintf("%s.%d.%s.%s", s[0], i, s[2], s[3])
		ips = genBClassIps(ip, num)
	}
	return ips
}

func GenIps(ip string, num []string, t string) [][]string {
	var ips [][]string
	if t == "b" {
		ips = genBClassIps(ip, num)
	}
	if t == "a" {
		ips = genAClassIps(ip, num)
	}
	if t == "172" {
		ips = gen172ClassIps(ip, num)
	}
	return ips
}
