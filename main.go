package main

import (
	"fmt"
	"strings"
	"sync"
	"zouyi/logagent/common"
	"zouyi/logagent/etcd"
	"zouyi/logagent/kafka"
	"zouyi/logagent/logger"
	"zouyi/logagent/setting"
	"zouyi/logagent/tailfile"
)

var (
	wg sync.WaitGroup
)

func main() {
	// 初始化日志设置
	logger.Init()

	var err error
	// 初始化配置信息
	if err = setting.Init(); err != nil {
		panic(fmt.Sprintf("init config failed, err:%v", err))
	}

	// 初始化kafka
	if err = kafka.Init(); err != nil {
		panic(fmt.Sprintf("init kafka failed, err:%v", err))
	}

	// 连接到etcd
	if err = etcd.Init(strings.Split(setting.Cfg.EtcdConfig.Address, ",")); err != nil {
		panic(fmt.Sprintf("init etcd failed, err:%v", err))
	}
	defer etcd.Close()

	// 获取ip
	if err = common.GetOutboundIP(); err != nil {
		panic(fmt.Sprintf("get local ip failed, err:%v\n", err))
	}
	// 获取etcd日志配置的key
	localLogKey := fmt.Sprintf(setting.Cfg.EtcdConfig.LogKey, common.LocalIp) // 根据每台服务器的主机来获取etcd中的日志path和topic
	// 获取日志并解析信息
	collectEntries, err := etcd.GetCollectEntries(localLogKey)
	if err != nil {
		panic(fmt.Sprintf("get conf from etcd err: %v", err))
	}
	// etcd监控日志信息变动
	//go etcd.WatchEntries(localLogKey)
	newEntryChan := etcd.WatchChan()

	// 初始化日志跟踪task
	err = tailfile.Init(collectEntries, newEntryChan)
	if err != nil {
		panic(fmt.Sprintf("init tailfile failed, err:%v", err))
	}

	// 启动
	run(localLogKey)
}

func run(logKey string) {
	for {

		wg.Add(1)
		go etcd.WatchEntries(logKey)
		// TODO, 监控本机信息
		wg.Wait()
	}
}
