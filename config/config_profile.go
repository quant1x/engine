package config

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"gitee.com/quant1x/gox/logger"
)

type PprofParameter struct {
	Enable bool `yaml:"enable" default:"false"` // 是否开启go tool pprof
	Port   int  `yaml:"port" default:"6060"`    // pprof web端口
}

// PprofEnable 获取配置中pprof开关
func PprofEnable() bool {
	return GlobalConfig.Runtime.Pprof.Enable
}

// StartPprof 启动性能分析工具
func StartPprof() {
	if !PprofEnable() {
		return
	}
	go func() {
		addr := fmt.Sprintf("localhost:%d", GlobalConfig.Runtime.Pprof.Port)
		err := http.ListenAndServe(addr, nil)
		logger.Info("启动pprof性能分析工具", err)
	}()
}
