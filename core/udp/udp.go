package udp

import (
	"github.com/urfave/cli/v2"
	"net"
	. "netspy/core/log"
	"netspy/core/spy"
	"strconv"
	"time"
)

var (
	timeout time.Duration
	ports   []int
)

func send(addr string) bool {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return false
	}

	conn, err := net.DialTimeout("udp", udpAddr.String(), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	packet := []byte{0xFF, 0xFF, 0xFF, 0xFF}
	// If UDP packet is set - send it
	if len(packet) > 0 {
		_, err = conn.Write(packet)
		if err != nil {
			return false
		}
	}

	// Wait for at least 1 byte response
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	data := make([]byte, 1)
	_, err = conn.Read(data)
	if err != nil {
		return false
	}
	return true
}

func check(ip string) bool {
	for _, port := range ports {
		netloc := net.JoinHostPort(ip, strconv.Itoa(port))
		if send(netloc) {
			Log.Debugf("%s open", netloc)
			return true
		}
	}
	return false
}

func Spy(c *cli.Context) {
	Log.Info("use udp protocol to spy")
	timeout = time.Duration(c.Int("timeout")) * time.Second
	ports = c.IntSlice("port")
	spy.Spy(c, check)
}
