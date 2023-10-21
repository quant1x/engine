package services

import (
	"errors"
	"fmt"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/cron"
	"gitee.com/quant1x/gox/logger"
	"sync"
)

// Job 定时任务
//
//	默认每10秒检测1次
//	排名不分先后
type Job struct {
	name    string
	spec    string
	Service func()
}

var (
	ErrAlreadyExists = errors.New("job is already exists")
)

var (
	jobMutex sync.Mutex
	mapJobs  = map[string]Job{}
	crontab  = cron.New()
)

// Register 注册任务
func Register(name, spec string, callback func()) error {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	_, ok := mapJobs[name]
	if ok {
		return ErrAlreadyExists
	}
	if len(spec) == 0 {
		spec = "@every 10s"
	}
	job := Job{name: name, spec: spec, Service: callback}
	mapJobs[job.name] = job
	return nil
}

// DaemonService 守护进程服务入口
func DaemonService() {
	// 启动服务
	logger.Infof("启动定时任务列表")
	crontab.Start()

	for _, v := range mapJobs {
		message := fmt.Sprintf("Service: %s, Interval: %s", v.name, v.spec)
		logger.Info(message)
		_, err := crontab.AddFuncWithSkipIfStillRunning(v.spec, v.Service)
		if err != nil {
			logger.Infof(message+"failed, err: %s", err.Error())
		} else {
			logger.Infof(message + "success")
		}
	}
	// 等待结束
	coroutine.WaitForShutdown()
	// 关闭任务调度
	crontab.Stop()
}

func List() {
	for _, v := range mapJobs {
		spec := v.spec
		message := fmt.Sprintf("Service: %s, Interval: %s", v.name, spec)
		fmt.Println(message)
	}
}
