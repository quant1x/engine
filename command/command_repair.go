package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets/base"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	cmder "github.com/spf13/cobra"
	"strings"
	"time"
)

// CmdRepair 补登历史数据
var CmdRepair = &cmder.Command{
	Use:     "repair",
	Example: Application + " repair --all",
	//Args:    args.MinimumNArgs(0),
	Args: func(cmd *cmder.Command, args []string) error {
		return nil
	},
	Short: "修复股市数据",
	Long:  `修复股市数据`,
	Run: func(cmd *cmder.Command, args []string) {
		beginDate := trading.FixTradeDate(flagStartDate.Value)
		endDate := cache.DefaultCanReadDate()
		if len(flagEndDate.Value) > 0 {
			endDate = trading.FixTradeDate(flagEndDate.Value)
		}
		dates := trading.TradeRange(beginDate, endDate)
		count := len(dates)
		fmt.Printf("修复数据: %s => %s"+strings.Repeat("\r\n", 2), dates[0], dates[count-1])
		base.UpdateTickStartDate(dates[0])
		if flagAll.Value {
			handleRepairAll(dates)
		} else if flagBaseData.Value {
			keywords := []string{}
			if flagTrans.Value {
				keywords = append(keywords, flagTrans.Name)
			}
			if len(keywords) == 0 {
				handleRepairDataSet(dates)
			} else {
				plugins := cache.PluginsWithName(cache.PluginMaskBaseData, keywords...)
				if len(plugins) == 0 {
					fmt.Printf("没有找到名字是[%s]的数据插件\n", strings.Join(keywords, ","))
				} else {
					handleRepairData(dates, plugins)
				}
			}
		} else if flagFeatures.Value {
			handleRepairFeatures(dates)
		}
	},
}

func init() {
	commandInit(CmdRepair, &flagAll)
	commandInit(CmdRepair, &flagBaseData)
	commandInit(CmdRepair, &flagFeatures)
	commandInit(CmdRepair, &flagStartDate)
	commandInit(CmdRepair, &flagEndDate)
	commandInit(CmdRepair, &flagTrans)
}

func handleRepairAll(dates []string) {
	moduleName := "修复全部数据"
	count := len(dates)
	fmt.Println()
	fmt.Println()
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		barIndex++
		storages.RepairFeatures(&barIndex, cacheDate, featureDate)
		_ = cacheDate
		_ = featureDate
		bar.Add(1)
	}
	logger.Info("任务执行完毕.", time.Now())
	fmt.Println()
}

func handleRepairDataSet(dates []string) {
	fmt.Println()
	moduleName := "补登数据集合"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		barIndex++
		storages.RepairBaseData(&barIndex, cacheDate, featureDate)
		bar.Add(1)
	}
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}

func handleRepairFeatures(dates []string) {
	moduleName := "补登特征数据"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		barIndex++
		storages.RepairFeatures(&barIndex, cacheDate, featureDate)
		bar.Add(1)
	}
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}

func handleRepairData(dates []string, plugins []cache.DataPlugin) {
	fmt.Println()
	//// 1. 获取全部注册的数据集插件
	//mask := cache.PluginMaskBaseData
	////dataSetList := flash.DataSetList()
	//plugins := cache.Plugins(mask)

	moduleName := "修复数据"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		//barIndex++
		storages.RepairData(barIndex+1, cacheDate, featureDate, plugins)
		bar.Add(1)
	}
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}
