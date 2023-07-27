package common

import (
	"fmt"
	"net"
	"strings"
)

// LogEntry 需要收集的日志信息
type LogEntry struct {
	Path string `json:"path"`
	//Module string `json:"module"`
	Topic string `json:"topic"`
}

type CollectSysInfoConfig struct {
	Interval int64  `json:"interval"`
	Topic    string `json:"topic"`
}

func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}
