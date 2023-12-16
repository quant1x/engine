package command

import (
	"fmt"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/gox/runtime"
	cmder "github.com/spf13/cobra"
	"strings"
)

var (
	strategyNumber = 0 // 策略编号
	businessDebug  = runtime.Debug()
)

var engineCmd = &cmder.Command{
	Use: Application,
	Run: func(cmd *cmder.Command, args []string) {
		model, err := models.CheckoutStrategy(strategyNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		barIndex := 1
		models.ExecuteStrategy(model, &barIndex)
	},
	PersistentPreRun: func(cmd *cmder.Command, args []string) {
		// 重置全局调试状态
		runtime.SetDebug(businessDebug)
	},
	PersistentPostRun: func(cmd *cmder.Command, args []string) {
		//
	},
}

// 初始化全部子命令
func initSubCommands() {
	initPrint()
	initRepair()
	initUpdate()
	initRules()
	initSafes()
	initBackTesting()
}

// GlobalFlags engine支持的全部命令
func GlobalFlags() *cmder.Command {
	initSubCommands()
	engineCmd.Flags().IntVar(&strategyNumber, "strategy", models.DefaultStrategy, models.UsageStrategyList())
	engineCmd.Flags().IntVar(&models.CountDays, "count", 0, "统计多少天")
	engineCmd.Flags().IntVar(&models.CountTopN, "top", models.AllStockTopN(), "输出前排几名")
	engineCmd.PersistentFlags().BoolVar(&businessDebug, "debug", businessDebug, "打开业务调试开关, 慎重使用!")
	engineCmd.AddCommand(CmdVersion, CmdPrint, CmdBackTesting, CmdRules)
	engineCmd.AddCommand(CmdUpdate, CmdRepair, CmdService, CmdSafes)
	return engineCmd
}

func parseFlagError(err error) (flag, value string) {
	before, _, ok := strings.Cut(err.Error(), "flag:")
	if !ok {
		return
	}
	before = strings.TrimSpace(before)
	//_, err1 := fmt.Sscanf(before, "invalid argument \"%s\" for \"--%s\"", &value, &flag)
	//if err1 != nil {
	//	return
	//}
	arr := strings.Split(before, "\"")
	if len(arr) != 5 {
		return
	}
	value = strings.TrimSpace(arr[1])
	flag = strings.TrimSpace(arr[3])
	arr = strings.Split(flag, "-")
	flag = arr[len(arr)-1]
	return
}
