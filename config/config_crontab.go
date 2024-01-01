package config

// JobParameter 定时任务配置
type JobParameter struct {
	//Name    string `yaml:"name" default:""`       // 任务名称
	Trigger string `yaml:"trigger"  default:""`   // 触发条件
	Enable  bool   `yaml:"enable" default:"true"` // 任务是否有效
}

// CrontabConfig 获取定时任务配置
func CrontabConfig() map[string]JobParameter {
	return GlobalConfig.Runtime.Crontab
}

// GetJobParameter 获取计划执行任务
func GetJobParameter(name string) *JobParameter {
	mapJob := CrontabConfig()
	if len(mapJob) == 0 {
		return nil
	}
	v, ok := mapJob[name]
	if ok {
		return &v
	}
	return nil
}
