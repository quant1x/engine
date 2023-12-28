package command

import (
	"gitee.com/quant1x/engine/models"
	"gitee.com/quant1x/engine/tracker"
	"gitee.com/quant1x/gotdx/proto"
	cmder "github.com/spf13/cobra"
	"strings"
)

var (
	days         int    // 统计多少天
	topN         int    // 统计前N
	strategyCode int    // 策略ID
	securityCode string // 证券代码
)

// CmdBackTesting 回测
var CmdBackTesting = &cmder.Command{
	Use:   "backtesting",
	Short: "回测",
	Run: func(cmd *cmder.Command, args []string) {
		securityCode = strings.TrimSpace(securityCode)
		securityCode = proto.CorrectSecurityCode(safesSecurityCode)
		if len(securityCode) > 0 {
			tracker.CheckStrategy(strategyCode, securityCode)
		} else {
			tracker.BackTesting(days, topN)
		}
	},
}

func initBackTesting() {
	CmdBackTesting.Flags().IntVar(&days, "count", 0, "统计多少天")
	CmdBackTesting.Flags().IntVar(&topN, "top", models.AllStockTopN(), "输出前排几名")
	CmdBackTesting.Flags().IntVar(&strategyCode, "strategy", 0, "策略ID")
	CmdBackTesting.Flags().StringVar(&securityCode, "code", "", "证券代码")
}
