package etcd

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"time"
	"zouyi/logagent/common"
)

var (
	cli         *clientv3.Client
	entriesChan chan []common.LogEntry
)

func Init(address []string) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		logrus.Errorf("connect to etcd client failed, %v\n", err)
		return
	}

	entriesChan = make(chan []common.LogEntry)
	logrus.Info("connect to etcd client success")
	return

}

func GetCollectEntries(key string) (collectEntries []common.LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := cli.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, %v\n", err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Error("get conf from etcd failed, len(kvs) == 0\n")
		return
	}

	keyValues := resp.Kvs[0]
	err = json.Unmarshal(keyValues.Value, &collectEntries)
	if err != nil {
		logrus.Errorf("json unmarshal conf failed, err: %v\n", err)
		return nil, err
	}
	logrus.Debugf("load conf from etcd success, collectEntries:%#v", collectEntries)
	return
}

func WatchEntries(key string) {
	for {
		watchChan := cli.Watch(context.Background(), key)

		for watchResponse := range watchChan {
			if err := watchResponse.Err(); err != nil {
				logrus.Warnf("watch key:%s err:%s\n", key, err)
				continue
			}

			for _, event := range watchResponse.Events {
				logrus.Debugf("type:%s key:%s value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
				var newCollectEntries []common.LogEntry
				if event.Type == clientv3.EventTypeDelete {
					entriesChan <- newCollectEntries
					continue
				}
				err := json.Unmarshal(event.Kv.Value, &newCollectEntries)
				if err != nil {
					logrus.Errorf("unmarshal conf failed, err:%s\n", err)
					continue
				}
				entriesChan <- newCollectEntries
				logrus.Debug("send newCollectEntries to entriesChan success")
			}
		}
	}

}

func WatchChan() <-chan []common.LogEntry {
	return entriesChan
}

func Close() (err error) {
	err = cli.Close()
	return
}
