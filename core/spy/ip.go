package spy

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func genRandNum(count int, exist []string) []string {
	remain := 256 - len(exist)
	if count >= remain {
		count = remain
	}
	var randNum []string
	for {
		n := rand.Intn(256)
		s := strconv.Itoa(n)
		_, isFind := find(exist, s)
		if isFind {
			continue
		}
		_, isFind = find(randNum, s)
		if isFind {
			continue
		}
		randNum = append(randNum, s)
		if len(randNum) == count {
			break
		}
	}
	return randNum
}

func genBClassIps(ip string, num []string) [][]string {
	var ips [][]string
	s := strings.Split(ip, ".")
	randNum := genRandNum(3, num)
	for i := 0; i < 255; i++ {
		var ipg []string
		for _, v := range num {
			// ip group
			ipg = append(ipg, fmt.Sprintf("%s.%s.%d.%s", s[0], s[1], i, v))
		}
		for _, v := range randNum {
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
