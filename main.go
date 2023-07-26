package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	logger "zouyi/logagent/Logger"
	"zouyi/logagent/common"
	"zouyi/logagent/kafka"
	"zouyi/logagent/tailfile"
)

var (
	log *logrus.Logger
	wg  sync.WaitGroup
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"chan_size"`
}

type CollectConfig struct {
	Logfile string `ini:"logfile"`
}

type EtcdConfig struct {
	Address           string `ini:"address"`
	CollectLogKey     string `ini:"collect_log_key"`
	CollectSysInfoKey string `ini:"collect_sysinfo_key"`
}

func run(logConfKey string, sysInfoConf *common.CollectSysInfoConfig) {

}

func main() {
	logger.Init()

	var cfg Config
	// 1. cfg
	err := ini.MapTo(&cfg, "./conf/config.ini")
	if err != nil {
		panic(fmt.Sprintf("init config failed, err:%v", err))
	}

	// 2. kafka
	err = kafka.Init(strings.Split(cfg.KafkaConfig.Address, ","), cfg.KafkaConfig.ChanSize)
	if err != nil {
		panic(fmt.Sprintf("init kafka failed, err:%v", err))
	}
	logrus.Info("init kafka success")
	// 3. etcd

	// 4. tail
	err = tailfile.Init(cfg.CollectConfig.Logfile)
	logrus.Info("init tailfile success")

}
