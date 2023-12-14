package main

import (
	"fmt"
	"gitee.com/quant1x/engine/command"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"
)

var (
	MinVersion = "1.0.0" // 版本号

)

var (
	cpuProfile = "./cpu.pprof"
	memProfile = "./mem.pprof"
)

// 更新日线数据工具
func main() {
	fCpu, err := os.Create(cpuProfile)
	if err != nil {
		logger.Fatal(err)
	}
	_ = pprof.StartCPUProfile(fCpu)
	defer pprof.StopCPUProfile()
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
	err = cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
	fMem, err := os.Create(memProfile)
	if err != nil {
		logger.Fatal(err)
	}
	_ = pprof.WriteHeapProfile(fMem)
	_ = fMem.Close()
}
