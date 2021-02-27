package scheduler

import (
	"context"
	"encoding/json"
	"gin-etcd-task/common"
	"gin-etcd-task/db"
	"gin-etcd-task/model"
	"github.com/coreos/etcd/clientv3"
	"github.com/robfig/cron"
	"log"
	"strings"
	"sync"
)

// 传入的事件类型
type PushTask struct {
	Task *model.Task
	Type int // 传入的类型，把他还是得Del

}

// 传入的并发
type TaskChan struct {
	PushTask chan *PushTask
	Mu       *sync.Mutex
}

// 运行的并发
type RunTaskChan struct {
	RunTask chan *model.Task
	Mu      *sync.Mutex
}

// 并发管理器
type TaskChanMange struct {
	TaskChan TaskChan // 这里有问题，需要看一下
	RunTaskChan RunTaskChan
	CronTable map[string]*cron.Cron
}

// 初始化并发管理器
func InitTaskChan() (t *TaskChanMange) {
	var pushMu sync.Mutex
	var runMu sync.Mutex
	t =  &TaskChanMange{
		TaskChan{PushTask: make(chan *PushTask,1000), Mu: &pushMu},
		RunTaskChan{RunTask: make(chan *model.Task,1000), Mu: &runMu},
		make(map[string]*cron.Cron,10),
	}
	return
}


func (t *TaskChanMange) PushTask(pt *PushTask) {
	t.TaskChan.Mu.Lock()
	go Exec(pt.Task.Bash) // 在释放提交之前，把他运行一遍
	defer t.TaskChan.Mu.Unlock()
}


func (t *TaskChanMange) WatchTaskChange() {
	resKv, err := db.Kv.Get(context.TODO(), common.ETCD_TASK_NAME, clientv3.WithPrefix())
	if err != nil {
		log.Println(err.Error())
	}
	go func() {
		rev := resKv.Header.Revision + 1
		watchChan := db.Watcher.Watch(
			context.TODO(),
			common.ETCD_TASK_NAME,
			clientv3.WithRev(rev),
			clientv3.WithPrefix(),
		)
		for response := range watchChan {
			for _, event := range response.Events {
				switch event.Type {
				case 0: // 这里是推送
					var mt model.Task
					err := json.Unmarshal(event.Kv.Value,&mt)
					if err != nil {
						log.Fatal(err)
					}
					p := PushTask{Type: 1, Task: &mt}
					go t.StartCron(&p)
				case 1: // 这是删除
					name := strings.TrimPrefix(string(event.Kv.Key), common.ETCD_TASK_NAME)
					t.EndCron(name)
				}
			}
		}
	}()
}


// 开始定时器，有问题需要传入PushTask
func (t *TaskChanMange) StartCron(task *PushTask) {
	c := cron.New()
	err := c.AddFunc(task.Task.Cron, func() {
		go t.PushTask(task)
	})

	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	t.CronTable[task.Task.Name] = c
}

// 结束corn
func (t *TaskChanMange) EndCron(name string) bool {
	_, err := db.Kv.Delete(context.TODO(), common.ETCD_TASK_NAME+name, clientv3.WithPrefix())
	_, err = db.Kv.Delete(context.TODO(), common.ETCD_RUN_TASK+name, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
		return false
	}
	// 如果是有
	if t.CronTable[name] != nil {
		t.CronTable[name].Stop()
		delete(t.CronTable, name)
		return true
	}
	return true
}
