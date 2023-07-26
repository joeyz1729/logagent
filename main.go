package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
	logger "zouyi/logagent/Logger"
	"zouyi/logagent/kafka"
	"zouyi/logagent/setting"
	"zouyi/logagent/tailfile"
)

//func run(logConfKey string, sysInfoConf *common.CollectSysInfoConfig) {
//
//}

func run() (err error) {
	for {
		line, ok := <-tailfile.TailObj.Lines
		if !ok {
			logrus.Warningf("tail file close reopen, filename: %s\n", tailfile.TailObj.Filename)
			time.Sleep(time.Second)
			continue
		}
		msg := &kafka.Message{
			Topic: "web_log",
			Data:  line.Text,
		}
		logrus.Info("msg: ", line.Text)
		kafka.MsgChan <- msg
	}
}

func main() {
	logger.Init()
	err := setting.Init()

	// 2. kafka
	err = kafka.Init(strings.Split(setting.Cfg.KafkaConfig.Address, ","), setting.Cfg.KafkaConfig.ChanSize)
	if err != nil {
		panic(fmt.Sprintf("init kafka failed, err:%v", err))
	}
	logrus.Info("init kafka success")
	// 3. etcd

	// 4. tail
	err = tailfile.Init(setting.Cfg.CollectConfig.Logfile)
	logrus.Info("init tailfile success")

	err = run()
	if err != nil {
		return
	}
}
