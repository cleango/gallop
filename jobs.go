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

//jobOption 脚本属性
type jobOption func(g *jobOptions)

type jobOptions struct {
	oneRun bool
}

func newJobOptions() *jobOptions {
	return &jobOptions{
		oneRun: true,
	}
}

type JobOption func(op *jobOptions)

//JobOneOption 对局属性
func JobOneOption(b bool) JobOption {
	return func(g *jobOptions) {
		g.oneRun = b
	}
}

type jobExecute struct {
	job     cron.Job
	isRun   bool
	options *jobOptions
	id      int
}

func newJobExecute(job cron.Job, options *jobOptions) *jobExecute {
	return &jobExecute{job: job, options: options}
}

func (j *jobExecute) Run() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err.(error))
		}
	}()
	if j.options.oneRun {
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
func AddJob(spec string, job cron.Job, options ...JobOption) (int, error) {
	opt := newJobOptions()
	for _, v := range options {
		v(opt)
	}
	ex := newJobExecute(job, opt)
	id, err := getCronTask().AddJob(spec, ex)
	ex.id = int(id)
	return int(id), err
}

//RemoveJob 删除脚本
func RemoveJob(id int) {
	getCronTask().Remove(cron.EntryID(id))
}

//Job 注入脚本支持依赖对象
func (g *Gallop) Job(spec string, job cron.Job, options ...JobOption) *Gallop {
	aop.Provide(&inject.Object{
		Value: job,
	})

	AddJob(spec, job)
	return g
}
