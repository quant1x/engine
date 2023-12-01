package config

// Job 定时任务配置
type Job struct {
	Name    string `yaml:"name" default:""`     // 任务名称
	Trigger string `yaml:"trigger"  default:""` // 触发条件
	Enable  bool   `yaml:"enable" default:"true"`
}
