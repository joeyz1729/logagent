package common

import (
	"net"
	"strings"
)

// LogEntry 需要收集的日志信息
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
	//Module string `json:"module"`
}

type CollectSysInfoConfig struct {
	Interval int64  `json:"interval"`
	Topic    string `json:"topic"`
}

var LocalIp string

func GetOutboundIP() (err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	//fmt.Println(localAddr.String())
	LocalIp = strings.Split(localAddr.IP.String(), ":")[0]
	return
}
