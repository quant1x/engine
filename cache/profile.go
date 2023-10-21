package cache

import (
	"fmt"
	"gitee.com/quant1x/gox/logger"
	"net/http"
	_ "net/http/pprof"
)

// 启动性能分析工具
func startPprof() {
	go func() {
		addr := fmt.Sprintf("localhost:%d", EngineConfig.Runtime.Pprof.Port)
		err := http.ListenAndServe(addr, nil)
		logger.Info("启动pprof性能分析工具", err)
	}()
}
