package tailfile

import (
	"context"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"time"
	"zouyi/logagent/common"
	"zouyi/logagent/kafka"
)

var (
	localIP string
	cfg     = tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
)

type LogData struct {
	IP   string `json:"ip"`
	Data string `json:"data"`
}

type Task struct {
	path string
	//module   string
	topic    string
	instance *tail.Tail
	ctx      context.Context
	cancel   context.CancelFunc
}

func Init(collectEntries []common.CollectEntry) (err error) {

	for _, ce := range collectEntries {
		task, err := NewTask(ce)
		if err != nil {
			logrus.Errorf("tailfile: create new task path:%s topic: %s failed, err: %v\n", ce.Path, ce.Topic, err)
			// TODO
			return err
		}
		go run(task)
	}
	logrus.Info("tailfile init tasks success")
	return nil
}

//func Init() {
//	logrus.Info("etcd:init log success")
//	var err error
//	localIP, err = common.GetOutboundIP()
//	if err != nil {
//		logrus.Errorf("get local ip failed, %v", err)
//	}
//}

func NewTask(ce common.CollectEntry) (task *Task, err error) {
	task = &Task{
		path:  ce.Path,
		topic: ce.Topic,
	}
	ins, err := tail.TailFile(ce.Path, cfg)
	if err != nil {
		return nil, err
	}
	task.instance = ins
	ctx, cancel := context.WithCancel(context.Background())
	task.ctx = ctx
	task.cancel = cancel
	return task, nil
}

func run(t *Task) (err error) {
	logrus.Debugf("begin task [path:%s topic:%s] stop\n", t.path, t.topic)
	for {
		select {
		case <-t.ctx.Done():
			logrus.Warnf("the task [path:%s topic:%s] stop", t.path, t.topic)
			t.instance.Cleanup()
			return
		case line, ok := <-t.instance.Lines:
			if !ok {
				logrus.Errorf("task [path:%s topic:%s] read line failed", t.path, t.topic)
				time.Sleep(time.Second)
				continue
			}
			if len(line.Text) == 0 {
				logrus.Warn("line does not contain data, skip")
				continue
			}
			// 读取到日志中创建的非空行，创建msg并发送给kafka
			msg := &kafka.Message{
				Topic: t.topic,
				Data:  line.Text,
			}
			logrus.Debug("msg: ", line.Text)
			if err = kafka.SendLog(msg); err != nil {
				logrus.Warning("msgChan is full")
				continue
			}
		}
		logrus.Debug("send msg to kafka success")
	}
}
