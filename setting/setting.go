package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

type Config struct {
	KafkaConfig `ini:"kafka"`
	//CollectConfig `ini:"collect"`
	EtcdConfig `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"chan_size"`
}

//type CollectConfig struct {
//	Logfile string `ini:"logfile"`
//	// TODO
//}

type EtcdConfig struct {
	Address           string `ini:"address"`
	CollectLogKey     string `ini:"collect_log_key"`
	CollectSysInfoKey string `ini:"collect_sysinfo_key"`
}

var Cfg Config

func Init() (err error) {
	err = ini.MapTo(&Cfg, "./conf/config.ini")
	if err != nil {
		panic(fmt.Sprintf("init config failed, err:%v", err))
	}
	logrus.Info("init config success")
	return
}
