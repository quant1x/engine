package config

import (
	"fmt"
	"gitee.com/quant1x/gox/logger"
	"net/http"
	_ "net/http/pprof"
)

type PprofParameter struct {
	Enable bool `yaml:"enable" default:"true"` // 是否开启go tool pprof
	Port   int  `yaml:"port" default:"6060"`   // pprof web端口
}

// StartPprof 启动性能分析工具
func StartPprof() {
	if !EngineConfig.Runtime.Pprof.Enable {
		return
	}
	go func() {
		addr := fmt.Sprintf("localhost:%d", EngineConfig.Runtime.Pprof.Port)
		err := http.ListenAndServe(addr, nil)
		logger.Info("启动pprof性能分析工具", err)
	}()
}
