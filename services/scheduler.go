package services

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/quant1x/engine/config"
	"github.com/quant1x/x/coroutine"
	"github.com/quant1x/x/cron"
	"github.com/quant1x/x/logger"
	"github.com/quant1x/x/runtime"
)

// Task 定时任务
//
//	默认每10秒检测1次
//	排名不分先后
type Task struct {
	name    string // 任务名称
	spec    string // 触发条件
	Service func() // 任务函数
}

var (
	ErrAlreadyExists = errors.New("the job already exists") // 任务已经存在
	ErrForbidden     = errors.New("the job was forbidden")  // 任务被禁止
)

var (
	jobMutex sync.Mutex
	mapJobs  = map[string]Task{}
	crontab  = cron.New()
)

// Register 注册定时任务
func Register(name, spec string, callback func()) error {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	_, ok := mapJobs[name]
	if ok {
		return ErrAlreadyExists
	}
	if len(spec) == 0 {
		spec = CronDefaultInterval
	}
	enable := true
	jobParam := config.GetJobParameter(name)
	if jobParam != nil {
		enable = jobParam.Enable
		trigger := strings.TrimSpace(jobParam.Trigger)
		if len(trigger) > 0 {
			spec = trigger
		}
	}
	if !enable {
		// 配置禁止任务, 不能返回错误
		return nil
	}
	job := Task{name: name, spec: spec, Service: callback}
	mapJobs[job.name] = job
	return nil
}

// DaemonService 守护进程服务入口
func DaemonService() {
	jobMutex.Lock()
	// 启动服务
	logger.Infof("启动定时任务列表")
	crontab.Start()
	// 输出2个空白行
	fmt.Printf("%s", strings.Repeat("\n", 2))
	for _, v := range mapJobs {
		message := fmt.Sprintf("Service: %s, Interval: %s, ", v.name, v.spec)
		logger.Info(message)
		_, err := crontab.AddJobWithSkipIfStillRunning(v.spec, v.Service)
		if err != nil {
			logger.Infof(message+"failed, err: %s", err.Error())
		} else {
			logger.Infof(message + "success")
		}
	}
	jobMutex.Unlock()
	// 等待结束
	coroutine.WaitForShutdown()
	// 关闭任务调度
	crontab.Stop()
}

// PrintJobList 输出定时任务列表
func PrintJobList() {
	for _, v := range mapJobs {
		spec := v.spec
		message := fmt.Sprintf("Service: %s, Interval: %s, method: %s", v.name, spec, runtime.FuncName(v.Service))
		fmt.Println(message)
	}
}
