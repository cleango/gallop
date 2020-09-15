package gallop

import (
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

//AddJob 注入脚本,不支持
func AddJob(spec string, job cron.Job)(int,error) {
	id, err := getCronTask().AddJob(spec, job)
	return int(id),err
}
//Job 注入脚本支持依赖对象
func (g *Gallop) Job(spec string, job cron.Job) *Gallop {
	aop.Provide(&inject.Object{
		Value:    job,
	})
	AddJob(spec,job)
	return g
}
