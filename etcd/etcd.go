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
	cli      *clientv3.Client
	confChan chan []*common.CollectEntry
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

	confChan = make(chan []*common.CollectEntry)
	logrus.Info("connect to etcd client success")
	return

}

func GetConf(key string) (conf []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := cli.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, %v\n", err)
		return
	}

	if len(resp.Kvs) == 0 {
		// TODO
		logrus.Error("get conf from etcd failed, len(kvs) == 0\n")
		return
	}

	keyValues := resp.Kvs[0]
	err = json.Unmarshal(keyValues.Value, &conf)
	if err != nil {
		logrus.Errorf("json unmarshal conf failed, err: %v\n", err)
		return nil, err
	}
	logrus.Debugf("load conf from etcd success, conf:%#v", conf)
	return
}

func Close() (err error) {
	err = cli.Close()
	return
}
