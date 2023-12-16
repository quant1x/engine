package main

import (
	"fmt"
	"gitee.com/quant1x/engine/command"
	"gitee.com/quant1x/engine/config"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"
)

var (
	MinVersion = "1.0.0" // 应用版本号
	//TdxVersion    = "1.0.0" // gotdx版本号
	//PandasVersion = "1.0.0" // pandas版本号
	//TaLibVersion  = "1.0.0" // ta-lib版本号
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
		runtime.CatchPanic()
		elapsedTime := time.Since(mainStart) / time.Millisecond
		fmt.Printf("\n总耗时: %.3fs\n", float64(elapsedTime)/1000)
	}()
	// stock模块内的更新版本号
	command.UpdateApplicationVersion(MinVersion)
	runtime.GoMaxProcs()

	// 命令字
	cmd := command.GlobalFlags()
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
	if config.PprofEnable() {
		fMem, err := os.Create(memProfile)
		if err != nil {
			logger.Fatal(err)
		}
		_ = pprof.WriteHeapProfile(fMem)
		_ = fMem.Close()
	}
}
