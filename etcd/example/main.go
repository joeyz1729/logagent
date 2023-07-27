package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"zouyi/logagent/etcd"
)

func main() {
	if err := etcd.Init([]string{"127.0.0.1:2379"}); err != nil {
		panic(fmt.Sprintf("init etcd failed, err:%v", err))
	}
	defer etcd.Close()
	logKey := `/log_agent/log_key/192.168.1.100`
	e1 := `{"path":"/tmp/log-agent/shopping.log","topic":"shopping"}`
	e2 := `{"path":"/tmp/log-agent/web.log","topic":"web"}`
	e3 := `{"path":"/tmp/log-agent/db.log","topic":"db"}`
	entries := []string{e1, e2, e3}
	value := fmt.Sprintf("[%s]", strings.Join(entries, ","))
	//value1 := fmt.Sprintf(`[%s]`, entry1)
	err := etcd.PutEntry(logKey, value)
	if err != nil {
		logrus.Errorf("put log entry err: %s\n", err)
	}
	val, err := etcd.GetEntry(logKey)
	if err != nil {
		fmt.Errorf("get log entry err: %s", err)
		return
	}
	fmt.Println(val)
}
