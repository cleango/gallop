package gallop

import (
	"github.com/cleango/gallop/logger"
	"github.com/cleango/gallop/third_plugins/inject"
	"github.com/robfig/cron/v3"
	"sync"
)

var onceCron sync.Once
var taskCron *cron.Cron //定时任务

func getCronTask() *cron.Cron {
	onceCron.Do(func() {
		taskCron = cron.New(cron.WithSeconds())
		taskCron.Start()
	})
	return taskCron
}

type jobE struct {
	job    cron.Job
	isRun  bool
	isOnly bool
}

func newJobE(job cron.Job, only bool) *jobE {
	return &jobE{job: job, isOnly: only}
}

func (j *jobE) Run() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err.(error))
		}
	}()
	if j.isOnly {
		if j.isRun {
			return
		}
		defer func() {
			j.isRun = false
		}()
		j.isRun = true
	}
	j.job.Run()
}

//AddJob 注入脚本,不支持
func AddJob(spec string, job *jobE) (int, error) {
	id, err := getCronTask().AddJob(spec, job)
	return int(id), err
}

//RemoveJob 删除脚本
func RemoveJob(id int) {
	getCronTask().Remove(cron.EntryID(id))
}

//Job 注入脚本支持依赖对象
func (g *Gallop) Job(spec string, job cron.Job, params ...bool) *Gallop {
	isOnly := true
	if len(params) == 1 {
		isOnly = params[0]
	}
	aop.Provide(&inject.Object{
		Value: job,
	})
	AddJob(spec, newJobE(job, isOnly))
	return g
}
