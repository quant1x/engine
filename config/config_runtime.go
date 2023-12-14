package config

// RuntimeParameter 运行时配置参数
type RuntimeParameter struct {
	Pprof   PprofParameter          `name:"性能分析" yaml:"pprof"`
	Debug   bool                    `name:"业务调试开关" yaml:"debug" default:"false"`
	Crontab map[string]JobParameter `name:"定时任务" yaml:"crontab" default:"{}"`
}
