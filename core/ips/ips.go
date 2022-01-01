package ips

import (
	"fmt"
	"strings"
)

func genBClassIps(ip string, n []string) []string {
	var gateways []string
	s := strings.Split(ip, ".")
	for i := 0; i < 255; i++ {
		for _, v := range n {
			gateways = append(gateways, fmt.Sprintf("%s.%s.%d.%s", s[0], s[1], i, v))
		}
	}
	return gateways
}

func genAClassIps(ip string, n []string) []string {
	var gateways []string
	s := strings.Split(ip, ".")
	for i := 0; i < 255; i++ {
		ip = fmt.Sprintf("%s.%d.%s.%s", s[0], i, s[2], s[3])
		gateways = genBClassIps(ip, n)
	}
	return gateways
}

func gen172ClassIps(ip string, n []string) []string {
	var gateways []string
	s := strings.Split(ip, ".")
	for i := 16; i < 32; i++ {
		ip = fmt.Sprintf("%s.%d.%s.%s", s[0], i, s[2], s[3])
		gateways = genBClassIps(ip, n)
	}
	return gateways
}

func GenIps(ip string, n []string, t string) []string {
	var gateways []string
	if t == "b" {
		gateways = genBClassIps(ip, n)
	}
	if t == "a" {
		gateways = genAClassIps(ip, n)
	}
	if t == "172" {
		gateways = gen172ClassIps(ip, n)
	}
	return gateways
}
