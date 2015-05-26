package qlog

import (
	"github.com/qiniu/log"
	"sync"
	"time"
)

var (
	GlbTaskRunner *QTaskRunner
)

type QTask struct {
	Bucket string
	Date   string
}

type QTaskRunner struct {
	tasks   []*QTask
	rwMutex *sync.RWMutex
}

func (this *QTaskRunner) Init() {
	this.tasks = make([]*QTask, 0)
	this.rwMutex = new(sync.RWMutex)
}

func (this *QTaskRunner) AddTask(bucket string, date string) {
	this.rwMutex.Lock()
	this.tasks = append(this.tasks, &QTask{bucket, date})
	this.rwMutex.Unlock()
}

//TODO 后面改成支持并发处理的
func (this *QTaskRunner) Scheduler() {
	for {
		var task *QTask = nil
		this.rwMutex.Lock()
		if len(this.tasks) > 0 {
			task = this.tasks[0]
			this.tasks = this.tasks[1:]
		}
		this.rwMutex.Unlock()
		if task != nil {
			go func() {
				WriteLogStatus(task.Bucket, task.Date, false)
				//read sync settings
				syncSettings, lErr := GetLogSyncSettings(task.Bucket)
				if lErr != nil {
					log.Error("error load sync settings for bucket,", task.Bucket)
					return
				}
				//sync
				logSync := QLogSync{
					Bucket:              syncSettings.Bucket,
					SaveBucket:          syncSettings.SaveBucket,
					SaveBucketDomain:    syncSettings.SaveBucketDomain,
					IsSaveBucketPrivate: syncSettings.IsSaveBucketPrivate,
					Mac:                 GlbConf.Mac(),
				}
				paths, syncErr := logSync.Sync(task.Date)
				if syncErr != nil {
					log.Error(syncErr.Error())
					return
				}
				//parse log content
				parseErr := ParseLogContent(task.Bucket, task.Date, paths)
				if parseErr != nil {
					log.Error(parseErr)
					return
				}
				WriteLogStatus(task.Bucket, task.Date, true)
			}()
		} else {
			<-time.After(5 * time.Second)
		}
	}
}
