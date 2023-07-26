package main

import (
	"fmt"
	"strings"
	logger "zouyi/logagent/Logger"
	"zouyi/logagent/etcd"
	"zouyi/logagent/kafka"
	"zouyi/logagent/setting"
	"zouyi/logagent/tailfile"
)

func main() {
	// 初始化日志设置
	logger.Init()

	// 初始化配置信息
	err := setting.Init()
	if err != nil {
		panic(fmt.Sprintf("init config failed, err:%v", err))
	}
	// 初始化kafka
	err = kafka.Init()
	if err != nil {
		panic(fmt.Sprintf("init kafka failed, err:%v", err))
	}

	// 连接到etcd
	err = etcd.Init(strings.Split(setting.Cfg.EtcdConfig.Address, ","))
	if err != nil {
		panic(fmt.Sprintf("init etcd failed, err:%v", err))
	}
	defer etcd.Close()

	// 从etcd中获取需要管理的日志信息
	//ip, err := common.GetOutboundIP()
	//if err != nil {
	//	panic(fmt.Sprintf("get local ip failed, err:%v\n", err))
	//}
	//collectLogKey := fmt.Sprintf(setting.Cfg.EtcdConfig.CollectLogKey, ip)
	// etcdctl put k1 '[{"path":"/tmp/log-agent/shopping.log","topic":"shopping"},{"path":"/tmp/log-agent/web.log","topic":"web"}]'
	collectLogKey := "k1" // 根据每台服务器的主机来获取etcd中的日志path和topic
	collectEntries, err := etcd.GetConf(collectLogKey)
	if err != nil {
		panic(fmt.Sprintf("get conf from etcd err: %v", err))
	}

	err = tailfile.Init(collectEntries)
	if err != nil {
		panic(fmt.Sprintf("init tailfile failed, err:%v", err))
	}
	run()
}

func run() {
	for {
		select {}
	}
}
