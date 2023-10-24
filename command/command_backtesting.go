package command

import (
	"gitee.com/quant1x/engine/strategies"
	"gitee.com/quant1x/engine/tracker"
	cmder "github.com/spf13/cobra"
)

var (
	days int // 统计多少天
	topN int // 统计前N
)

// CmdBackTesting 回测
var CmdBackTesting = &cmder.Command{
	Use:   "backtesting",
	Short: "回测",
	Run: func(cmd *cmder.Command, args []string) {
		tracker.BackTesting(days, topN)
	},
}

func initBackTesting() {
	CmdBackTesting.Flags().IntVar(&days, "count", 0, "统计多少天")
	CmdBackTesting.Flags().IntVar(&topN, "top", strategies.AllStockTopN(), "输出前排几名")
}
