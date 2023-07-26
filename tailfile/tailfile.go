package tailfile

import (
	"context"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var (
	localIP string
	obj     *tail.Tail
)

type LogData struct {
	IP   string `json:"ip"`
	Data string `json:"data"`
}

type tailObj struct {
	path     string
	module   string
	topic    string
	instance *tail.Tail
	ctx      context.Context
	cancel   context.CancelFunc
}

func Init(filename string) (err error) {
	cfg := tail.Config{
		ReOpen:   true,
		Follow:   true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		Poll:     true,
	}
	obj, err = tail.TailFile(filename, cfg)
	if err != nil {
		logrus.Error("tailfile: create tailObj for path:%s failed, err: %v\n", filename, err)
		return
	}
	return
}

//func Init() {
//	logrus.Info("etcd:init log success")
//	var err error
//	localIP, err = common.GetOutboundIP()
//	if err != nil {
//		logrus.Errorf("get local ip failed, %v", err)
//	}
//}

//
//func NewTailObj(path, module, topic string) (tObj *tailObj, err error) {
//	tObj = &tailObj{
//		path:   path,
//		module: module,
//		topic:  topic,
//	}
//	ctx, cancel := context.WithCancel(context.Background())
//	tObj.ctx = ctx
//	tObj.cancel = cancel
//	err = tObj.Init()
//	return
//}
//
//func (t *tailObj) Init() (err error) {
//	t.instance, err = tail.TailFile(t.path, tail.Config{
//		ReOpen:   true,
//		Follow:   true,
//		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
//		Poll:     true,
//	})
//	if err != nil {
//		fmt.Println("init tail failed, err:", err)
//		return
//	}
//	return
//}
//
//func (t *tailObj) run() {
//	for {
//		select {
//		case <-t.ctx.Done():
//			logrus.Warnf("the task of path:%s is stop...", t.path)
//			t.instance.Cleanup()
//			return
//		case line, ok := <-t.instance.Lines:
//			if !ok {
//				logrus.Errorf("read line failed")
//				continue
//			}
//			data := &LogData{
//				IP:   localIP,
//				Data: line.Text,
//			}
//			jsonData, err := json.Marshal(data)
//			if err != nil {
//				logrus.Warningf("marshal tailfile.LogData failled, err:%v", err)
//			}
//			msg := &kafka.Message{
//				Data:  string(jsonData),
//				Topic: t.topic,
//			}
//			err = kafka.SendLog(msg)
//			if err != nil {
//				logrus.Errorf("send to kafka failed, err:%v\n", err)
//			}
//		}
//		logrus.Info("send msg to kafka success")
//	}
//}
