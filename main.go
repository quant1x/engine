package main

import (
	"fmt"
	"gitee.com/quant1x/engine/command"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/strategies"
	"gitee.com/quant1x/engine/tracker"
	"gitee.com/quant1x/engine/util/runtime"
	cmder "github.com/spf13/cobra"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"
)

var (
	MinVersion     = "1.0.0" // 版本号
	strategyNumber = 0       // 策略编号
)

var (
	cpuProfile = "./cpu.pprof"
	memProfile = "./mem.pprof"
)

// 更新日线数据工具
func main() {
	fCpu, err := os.Create(cpuProfile)
	if err != nil {
		log.Fatal(err)
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

	var rootCmd = &cmder.Command{
		Use: command.Application,
		Run: func(cmd *cmder.Command, args []string) {
			//stat.SetAvx2Enabled(modules.CpuAvx2)
			//runtime.GoMaxProcs(modules.CpuNum)
			var model models.Strategy
			switch strategyNumber {
			default:
				model = new(strategies.ModelNo1)
			}
			fmt.Printf("策略模块: %s\n", model.Name())
			if strategies.CountDays > 0 {
				tracker.BackTesting(strategies.CountDays, strategies.CountTopN)
			} else {
				// 执行策略
				barIndex := 1
				models.ExecuteStrategy(model, &barIndex)
			}
		},
	}
	rootCmd.Flags().IntVar(&strategyNumber, "strategy", models.DefaultStrategy, "策略编号")
	rootCmd.Flags().IntVar(&strategies.CountDays, "count", 0, "统计多少天")
	rootCmd.Flags().IntVar(&strategies.CountTopN, "top", strategies.AllStockTopN(), "输出前排几名")
	command.Init()
	rootCmd.AddCommand(command.CmdVersion, command.CmdPrint, command.CmdBackTesting, command.CmdRule)
	rootCmd.AddCommand(command.CmdUpdate, command.CmdRepair, command.CmdService)
	_ = rootCmd.Execute()
	fMem, err := os.Create(memProfile)
	if err != nil {
		log.Fatal(err)
	}
	_ = pprof.WriteHeapProfile(fMem)
	_ = fMem.Close()
}
