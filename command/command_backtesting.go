package command

import (
	"gitee.com/quant1x/engine/global"
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/tracker"
	"gitee.com/quant1x/exchange"
	cmder "github.com/spf13/cobra"
	"strings"
)

var (
	strategyCode uint64 // 策略ID
	securityCode string // 证券代码
	topN         int    // 统计前N
	days         int    // 统计多少天
	date         string // 回测日期
)

// CmdBackTesting 回测
var (
	CmdBackTesting *cmder.Command
)

func initBackTesting(variables *global.Variables) {
	CmdBackTesting = &cmder.Command{
		Use:   "backtesting",
		Short: "回测",
		Run: func(cmd *cmder.Command, args []string) {
			securityCode = strings.TrimSpace(securityCode)
			securityCode = exchange.CorrectSecurityCode(securityCode)
			if len(securityCode) > 0 {
				tracker.CheckStrategy(strategyCode, securityCode, date, *variables.MarketData)
			} else {
				tracker.BackTesting(strategyCode, days, topN, *variables.MarketData)
			}
		},
	}

	CmdBackTesting.Flags().IntVar(&days, "count", 0, "统计多少天")
	CmdBackTesting.Flags().IntVar(&topN, "top", models.AllStockTopN(), "输出前排几名")
	CmdBackTesting.Flags().Uint64Var(&strategyCode, "strategy", 0, "策略ID")
	CmdBackTesting.Flags().StringVar(&securityCode, "code", "", "证券代码")
	CmdBackTesting.Flags().StringVar(&date, "date", "", "日期")
}
