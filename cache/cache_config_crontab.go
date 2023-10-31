package cache

// Job 定时任务配置
type Job struct {
	Name    string `yaml:"name"`    // 任务名称
	Trigger string `yaml:"trigger"` // 触发条件
}
