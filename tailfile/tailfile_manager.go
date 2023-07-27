package tailfile

import (
	"github.com/sirupsen/logrus"
	"zouyi/logagent/common"
)

type taskManager struct {
	Tasks map[string]*Task
	//EntryList []common.LogEntry
	EntryChan <-chan []common.LogEntry
}

var Manager taskManager

func Init(collectEntries []common.LogEntry, newWatchChan <-chan []common.LogEntry) (err error) {
	Manager = taskManager{
		Tasks: make(map[string]*Task, 32),
		//EntryList: collectEntries,
		EntryChan: newWatchChan,
	}

	for _, ce := range collectEntries {
		task, err := NewTask(ce)
		if err != nil {
			logrus.Errorf("tailfile: create new task path:%s failed, err: %v\n", ce.Path, err)
			// TODO, 如果其中部分成功，部分失败
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
		logrus.Warn("changes in log entry")

		// 日志条目更新
		visited := make(map[string]struct{})
		for _, newEntry := range newEntries {
			if m.exist(newEntry) {
				visited[newEntry.Path+newEntry.Topic] = struct{}{}
				logrus.Debugf("task exists [path: %s]\n", newEntry.Path)
				continue
			}
			// 新的日志条目
			newTask, err := NewTask(newEntry)
			if err != nil {
				logrus.Errorf("create task [path: %s] err: %s\n", newEntry.Path, err)
				// TODO, 部分task创建失败时retry
				continue
			}
			visited[newEntry.Path+newEntry.Topic] = struct{}{}
			go newTask.run()
			Manager.Tasks[newEntry.Path+newEntry.Topic] = newTask
		}

		// 删除不需要的日志监控
		for key := range m.Tasks {
			// TODO, 优化集合操作
			if _, ok := visited[key]; ok {
				// 保留
				continue
			}
			m.Tasks[key].cancel()
			m.Tasks[key].instance.Stop()
			logrus.Debugf("task cancel [key:%s]", key)
			delete(m.Tasks, key)
		}
	}

}

func (m taskManager) exist(entry common.LogEntry) bool {
	if _, ok := m.Tasks[entry.Path+entry.Topic]; ok {
		// TODO key是采用path还是path+topic
		return true
	}
	return false
}
