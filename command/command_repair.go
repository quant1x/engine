package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	flags "github.com/spf13/cobra"
	"time"
)

// CmdRepair 补登历史数据
var CmdRepair = &flags.Command{
	Use:     "repair",
	Example: Application + " repair --all",
	//Args:    args.MinimumNArgs(0),
	Args: func(cmd *flags.Command, args []string) error {
		return nil
	},
	Short: "回补股市数据",
	Long:  `回补股市数据`,
	Run: func(cmd *flags.Command, args []string) {
		beginDate := trading.FixTradeDate(flagStartDate.Value)
		endDate := cachel5.DefaultCanReadDate()
		if len(flagEndDate.Value) > 0 {
			endDate = trading.FixTradeDate(flagEndDate.Value)
		}
		if flagHistory.Value {
			//date := "2023-09-28"
			//cacheDate, featureDate := cachel5.CorrectDate(date)
			//update.Repair(cacheDate, featureDate)
			handleRepair(beginDate, endDate)
		}
	},
}

func init() {
	flagHistory.init(CmdRepair)
	flagStartDate.init(CmdRepair)
	flagEndDate.init(CmdRepair)
}

func handleRepair(begin, end string) {
	moduleName := "补登历史数据"
	dates := trading.TradeRange(begin, end)
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		bar.Add(1)
		cacheDate, featureDate := cachel5.CorrectDate(date)
		storages.Repair(cacheDate, featureDate)
	}
	logger.Info("任务执行完毕.", time.Now())
	fmt.Println()
}
