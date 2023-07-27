package tailfile

import (
	"github.com/sirupsen/logrus"
	"zouyi/logagent/common"
)

type taskManager struct {
	Tasks     map[string]*Task
	EntryList []common.LogEntry
	EntryChan <-chan []common.LogEntry
}

var Manager taskManager

func Init(collectEntries []common.LogEntry, newWatchChan <-chan []common.LogEntry) (err error) {
	Manager = taskManager{
		Tasks:     make(map[string]*Task, 32),
		EntryList: collectEntries,
		EntryChan: newWatchChan,
	}

	for _, ce := range collectEntries {
		task, err := NewTask(ce)
		if err != nil {
			logrus.Errorf("tailfile: create new task path:%s topic: %s failed, err: %v\n", ce.Path, ce.Topic, err)
			// TODO
			return err
		}
		Manager.Tasks[task.path+task.topic] = task
		go task.run()
	}
	logrus.Info("tailfile init tasks success")

	// 监控etcd中日志配置的修改
	go Manager.watch()

	return nil
}

func (m taskManager) watch() {
	logrus.Debug("wait for new log entries")
	for {
		newEntries := <-m.EntryChan
		logrus.Debug("new entries")

		// 添加新的日志监控
		for _, newEntry := range newEntries {
			if m.exist(newEntry) {
				logrus.Debugf("task [path: %s, topic: %s] exists\n", newEntry.Path, newEntry.Topic)
				continue
			}
			// 新的日志条目
			newTask, err := NewTask(newEntry)
			if err != nil {
				logrus.Errorf("create task [path: %s, topic: %s] err: %s\n", newEntry.Path, newEntry.Topic, err)
				continue
			}
			go newTask.run()
			Manager.Tasks[newEntry.Path+newEntry.Topic] = newTask
		}

		// 删除不需要的日志监控
		for key := range m.Tasks {
			ok := false
			for _, newEntry := range newEntries {
				if key == newEntry.Path+newEntry.Topic {
					ok = true
					break
				}
			}
			if !ok {
				m.Tasks[key].cancel()
				logrus.Debugf("task [key:%s] cancel", key)
				delete(m.Tasks, key)
			}
		}
	}

}

func (m taskManager) exist(entry common.LogEntry) bool {
	if _, ok := m.Tasks[entry.Path+entry.Topic]; ok {
		// TODO 日志属性修改时，例如修改日志topic
		return true
	}
	return false
}
