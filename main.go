package main

import (
	"fmt"
	"gitee.com/quant1x/engine/command"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/util/functions"
	"gitee.com/quant1x/gox/logger"
	cmder "github.com/spf13/cobra"
	"runtime/debug"
	"time"
)

var (
	MinVersion     = "1.0.0" // 版本号
	strategyNumber = 0       // 策略编号
)

// 更新日线数据工具
func main() {
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
	rootCmd.AddCommand(command.CmdVersion, command.CmdPrint, command.CmdUpdate, command.CmdRepair)
	_ = rootCmd.Execute()
}
