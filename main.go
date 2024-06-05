package main

import (
	"fmt"
	"gitee.com/quant1x/engine/command"
	"gitee.com/quant1x/engine/config"
	_ "gitee.com/quant1x/engine/strategies"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"
)

// 系统构建时传入的值
// go build -ldflags "-X 'main.MinVersion=${version}'"

var (
	MinVersion = "1.0.0" // 应用版本号
)

var (
	cpuProfile = "./cpu.pprof"
	memProfile = "./mem.pprof"
)

// 更新日线数据工具
func main() {
	if config.PprofEnable() {
		fCpu, err := os.Create(cpuProfile)
		if err != nil {
			logger.Fatal(err)
		}
		_ = pprof.StartCPUProfile(fCpu)
		defer pprof.StopCPUProfile()
	}
	mainStart := time.Now()
	defer func() {
		runtime.CatchPanic("")
		elapsedTime := time.Since(mainStart) / time.Millisecond
		fmt.Printf("\n总耗时: %.3fs\n", float64(elapsedTime)/1000)
	}()
	// stock模块内的更新版本号
	command.UpdateApplicationVersion(MinVersion)
	runtime.GoMaxProcs()

	// 命令字
	rootCommand := command.GlobalFlags()
	_ = rootCommand.Execute()
	if config.PprofEnable() {
		fMem, err := os.Create(memProfile)
		if err != nil {
			logger.Fatal(err)
		}
		_ = pprof.WriteHeapProfile(fMem)
		_ = fMem.Close()
	}
}
