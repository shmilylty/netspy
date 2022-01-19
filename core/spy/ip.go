package spy

import (
	"math/rand"
	"net"
)

func find(slice []int, i int) (int, bool) {
	for index, value := range slice {
		if value == i {
			return index, true
		}
	}
	return -1, false
}

func genRandNum(count int, exist []int) []int {
	var randNum []int
	if count == 0 {
		return randNum
	}
	remain := 255 - len(exist)
	if count >= remain {
		count = remain
	}
	for {
		i := rand.Intn(256)
		if i == 0 {
			continue
		}
		_, isFind := find(exist, i)
		if isFind {
			continue
		}
		_, isFind = find(randNum, i)
		if isFind {
			continue
		}
		randNum = append(randNum, i)
		if len(randNum) == count {
			break
		}
	}
	return randNum
}

func GenIPS(netips []net.IP, endNum []int, count int) [][]string {
	var ips [][]string
	for _, ip := range netips {
		// ip group
		var ipg []string
		for _, i := range endNum {
			ip[3] = byte(i)
			ipg = append(ipg, ip.String())
		}
		// 每个段都随机生成IP尾数
		randNum := genRandNum(count, endNum)
		for _, i := range randNum {
			ip[3] = byte(i)
			ipg = append(ipg, ip.String())
		}
		ips = append(ips, ipg)
	}
	return ips
}
