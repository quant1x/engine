package main

import (
	"fmt"
	"gitee.com/quant1x/engine/command"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/util/functions"
	"gitee.com/quant1x/gox/logger"
	cmder "github.com/spf13/cobra"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/debug"
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
	//go func() {
	//	log.Println(http.ListenAndServe(":8000", nil))
	//}()
	mainStart := time.Now()
	defer func() {
		if err := recover(); err != nil {
			s := string(debug.Stack())
			fmt.Printf("\nerr=%v, stack=%s\n", err, s)
			logger.Fatalf("%s 异常: %+v", command.Application, err)
		}
		elapsedTime := time.Since(mainStart) / time.Millisecond
		fmt.Printf("\n总耗时: %.3fs\n", float64(elapsedTime)/1000)
	}()
	// stock模块内的更新版本号
	command.UpdateApplicationVersion(MinVersion)
	functions.GOMAXPROCS()

	var rootCmd = &cmder.Command{
		Use: command.Application,
		Run: func(cmd *cmder.Command, args []string) {
			//stat.SetAvx2Enabled(modules.CpuAvx2)
			//runtime.GOMAXPROCS(modules.CpuNum)
			var model models.Strategy
			switch strategyNumber {
			default:
				model = new(models.ModelNo1)
			}
			fmt.Printf("策略模块: %s\n", model.Name())
			// 执行策略
			barIndex := 1
			models.ExecuteStrategy(model, &barIndex)
		},
	}
	rootCmd.Flags().IntVar(&strategyNumber, "strategy", models.DefaultStrategy, "策略编号")
	command.Init()
	rootCmd.AddCommand(command.CmdVersion, command.CmdPrint, command.CmdUpdate, command.CmdRepair, command.CmdService)
	_ = rootCmd.Execute()
	fMem, err := os.Create(memProfile)
	if err != nil {
		log.Fatal(err)
	}
	_ = pprof.WriteHeapProfile(fMem)
	_ = fMem.Close()
}
