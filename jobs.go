package gallop

import (
	"github.com/cleango/gallop/third_plugins/inject"
	"github.com/robfig/cron/v3"
	"log"
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

//AddJob 注入脚本
func AddJob(spec string, job cron.Job) {
	aop.Provide(&inject.Object{
		Value:    job,
	})
	_, err := getCronTask().AddJob(spec, job)
	if err != nil {
		log.Println(err)
	}
}
//Job 注入脚本
func (g *Gallop) Job(spec string, job cron.Job) *Gallop {
	AddJob(spec,job)
	return g
}
