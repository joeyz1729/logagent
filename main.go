package main

import (
	"fmt"
	"strings"
	logger "zouyi/logagent/Logger"
	"zouyi/logagent/common"
	"zouyi/logagent/etcd"
	"zouyi/logagent/kafka"
	"zouyi/logagent/setting"
	"zouyi/logagent/tailfile"
)

func main() {
	// 初始化日志设置
	logger.Init()

	var err error
	// 初始化配置信息
	err = setting.Init()
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

	// 获取ip
	ip, err := common.GetOutboundIP()
	if err != nil {
		panic(fmt.Sprintf("get local ip failed, err:%v\n", err))
	}
	// 获取etcd日志配置的key
	logKey := fmt.Sprintf(setting.Cfg.EtcdConfig.LogKey, ip) // 根据每台服务器的主机来获取etcd中的日志path和topic
	// 获取日志信息
	collectEntries, err := etcd.GetCollectEntries(logKey)
	if err != nil {
		panic(fmt.Sprintf("get conf from etcd err: %v", err))
	}
	// etcd监控日志信息变动
	newEntryChan := etcd.WatchChan()

	// 初始化日志跟踪task
	err = tailfile.Init(collectEntries, newEntryChan)
	if err != nil {
		panic(fmt.Sprintf("init tailfile failed, err:%v", err))
	}

	// 启动
	run()
}

func run() {
	for {
		select {}
	}
}
