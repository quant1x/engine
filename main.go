package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"

	"gitee.com/quant1x/engine/command"
	"gitee.com/quant1x/engine/config"
	_ "gitee.com/quant1x/engine/strategies"
	"gitee.com/quant1x/engine/utils"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
)

// 系统构建时传入的值
// go build -ldflags "-X 'main.MinVersion=${version}'"

var (
	MinVersion = utils.InvalidVersion // 应用版本号
	tdxVersion = utils.InvalidVersion // tdx api版本号
)

var (
	cpuProfile = "./cpu.pprof"
	memProfile = "./mem.pprof"
)

func resetVersions() {
	if MinVersion == utils.InvalidVersion {
		MinVersion = utils.CurrentVersion()
	}
	if tdxVersion == utils.InvalidVersion {
		tdxVersion = utils.RequireVersion("gitee.com/quant1x/data/level1")
	}
}

// 更新基础数据,特征,执行策略,回测等功能入口
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
	resetVersions()
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
