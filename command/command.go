package command

import (
	"fmt"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/tracker"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/runtime"
	"gitee.com/quant1x/num"
	cmder "github.com/spf13/cobra"
	"os"
	goruntime "runtime"
	"strings"
	_ "unsafe" // For go:linkname
)

func cpuModelName() string {
	return "Unknown"
}

const (
	KB = 1024      // Kilo Byte
	MB = 1024 * KB // Mega Byte
	GB = 1024 * MB // Giga Byte
)

var (
	strategyNumber = uint64(1)              // 策略编号
	businessDebug  = runtime.Debug()        // 业务调试开关
	cpuAvx2        = false                  // AVX2加速状态
	cpuNum         = goruntime.NumCPU() / 2 // cpu数量
)

var (
	commandDefaultLongFlag = "" // 默认的长标志为空, 主要用在tools
)

// 输出欢迎语
func printMotd() {
	fmt.Printf("CPU: %s, %dCores, AVX2: %t\n", cpuModelName(), goruntime.NumCPU(), num.GetAvx2Enabled())
	fmt.Println()
}

// 初始化全部子命令
func initSubCommands() {
	initPrint()
	initRepair()
	initUpdate()
	initRules()
	initSafes()
	initBackTesting()
	initTracker()
	initTools()
	initService()
}

// InitCommands 公开初始化函数
func InitCommands() {
	initSubCommands()
}

// GlobalFlags engine支持的全部命令
func GlobalFlags() *cmder.Command {
	initSubCommands()
	engineCmd := &cmder.Command{
		Use: Application,
		Run: func(cmd *cmder.Command, args []string) {
			logger.Warnf("stock default args:%+v", os.Args)
			model, err := models.CheckoutStrategy(strategyNumber)
			if err != nil {
				fmt.Println(err)
				return
			}
			// 输出欢迎语
			printMotd()
			barIndex := 1
			tracker.ExecuteStrategy(model, &barIndex)
		},
		PersistentPreRun: func(cmd *cmder.Command, args []string) {
			// 重置全局调试状态
			runtime.SetDebug(businessDebug)
			// AVX2 加速
			num.SetAvx2Enabled(cpuAvx2)
			// 设置CPU最大核数
			runtime.GoMaxProcs(cpuNum)
		},
		PersistentPostRun: func(cmd *cmder.Command, args []string) {
			//
		},
	}
	engineCmd.Flags().Uint64Var(&strategyNumber, "strategy", models.DefaultStrategy, models.UsageStrategyList())
	engineCmd.Flags().IntVar(&models.CountDays, "count", 0, "统计多少天")
	engineCmd.Flags().IntVar(&models.CountTopN, "top", models.AllStockTopN(), "输出前排几名")
	engineCmd.PersistentFlags().BoolVar(&businessDebug, "debug", businessDebug, "打开业务调试开关, 慎重使用!")
	engineCmd.PersistentFlags().BoolVar(&cpuAvx2, "avx2", false, "Avx2 加速开关")
	engineCmd.PersistentFlags().IntVar(&cpuNum, "cpu", cpuNum, "设置CPU最大核数")
	engineCmd.AddCommand(CmdVersion, CmdSafes, CmdBestIP, CmdConfig, CmdTools)
	engineCmd.AddCommand(CmdUpdate, CmdRepair, CmdPrint)
	engineCmd.AddCommand(CmdBackTesting, CmdRules, CmdTracker)
	engineCmd.AddCommand(CmdService)
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
