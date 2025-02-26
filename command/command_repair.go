package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/storages"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	cmder "github.com/spf13/cobra"
	"strings"
	"time"
)

const (
	repairCommand     = "repair"
	repairDescription = "修复数据"
)

var (
	// CmdRepair 补登历史数据
	CmdRepair *cmder.Command = nil
)

func initRepair() {
	CmdRepair = &cmder.Command{
		Use:     repairCommand,
		Example: Application + " " + repairCommand + " --all",
		//Args:    args.MinimumNArgs(0),
		Args: func(cmd *cmder.Command, args []string) error {
			return nil
		},
		Short: repairDescription,
		Long:  repairDescription,
		Run: func(cmd *cmder.Command, args []string) {
			beginDate := exchange.FixTradeDate(flagStartDate.Value)
			endDate := cache.DefaultCanReadDate()
			if len(flagEndDate.Value) > 0 {
				endDate = exchange.FixTradeDate(flagEndDate.Value)
			}
			dates := exchange.TradingDateRange(beginDate, endDate)
			count := len(dates)
			if count == 0 {
				fmt.Printf("start=%s ~ end=%s 休市, 没有数据\n", beginDate, endDate)
				return
			}
			fmt.Printf("修复数据: %s => %s"+strings.Repeat("\r\n", 2), dates[0], dates[count-1])
			base.UpdateBeginDateOfHistoricalTradingData(dates[0])
			if flagAll.Value {
				handleRepairAll(dates)
			} else if len(flagBaseData.Value) > 0 {
				all, keywords := parseFields(flagBaseData.Value)
				if all || len(keywords) == 0 {
					handleRepairAllDataSets(dates)
				} else {
					plugins := cache.PluginsWithName(cache.PluginMaskBaseData, keywords...)
					if len(plugins) == 0 {
						fmt.Printf("没有找到名字是[%s]的数据插件\n", strings.Join(keywords, ","))
					} else {
						handleRepairDataSetsWithPlugins(dates, plugins)
					}
				}
			} else if len(flagFeatures.Value) > 0 {
				all, keywords := parseFields(flagFeatures.Value)
				if all || len(keywords) == 0 {
					handleRepairAllFeatures(dates)
				} else {
					plugins := cache.PluginsWithName(cache.PluginMaskFeature, keywords...)
					if len(plugins) == 0 {
						fmt.Printf("没有找到名字是[%s]的数据插件\n", strings.Join(keywords, ","))
					} else {
						handleRepairFeaturesWithPlugins(dates, plugins)
					}
				}
			} else {
				fmt.Println("Error: 非全部修复, 必须携带--features或--base")
				_ = cmd.Usage()
			}

		},
	}
	commandInit(CmdRepair, &flagAll)
	commandInit(CmdRepair, &flagStartDate)
	commandInit(CmdRepair, &flagEndDate)

	// 1. 基础数据
	plugins := cache.Plugins(cache.PluginMaskBaseData)
	flagBaseData.Usage = getPluginsUsage(plugins)
	commandInit(CmdRepair, &flagBaseData)

	// 2. 特征数据
	plugins = cache.Plugins(cache.PluginMaskFeature)
	flagFeatures.Usage = getPluginsUsage(plugins)
	commandInit(CmdRepair, &flagFeatures)
}

func handleRepairAll(dates []string) {
	handleRepairAllDataSets(dates)
	handleRepairAllFeatures(dates)
}

func handleRepairAllDataSets(dates []string) {
	moduleName := "补登数据集合"
	logger.Info(moduleName + ", 任务开始")
	mask := cache.PluginMaskBaseData
	plugins := cache.Plugins(mask)
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	barIndex++
	for _, date := range dates {
		//cacheDate, featureDate := cache.CorrectDate(date)
		storages.BaseDataUpdate(barIndex, date, plugins, cache.OpRepair)
		bar.Add(1)
	}
	bar.Wait()
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}

// 修复 - 指定的基础数据
func handleRepairDataSetsWithPlugins(dates []string, plugins []cache.DataAdapter) {
	moduleName := "修复数据"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		//cacheDate, featureDate := cache.CorrectDate(date)
		//barIndex++
		storages.BaseDataUpdate(barIndex+1, date, plugins, cache.OpRepair)
		bar.Add(1)
	}
	bar.Wait()
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}

// 修复 - 特征数据
func handleRepairAllFeatures(dates []string) {
	moduleName := "补登特征数据"
	logger.Info(moduleName + ", 任务开始")
	mask := cache.PluginMaskFeature
	plugins := cache.Plugins(mask)
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	barIndex++
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpRepair)
		bar.Add(1)
	}
	bar.Wait()
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}

// 修复 - 指定的特征数据
func handleRepairFeaturesWithPlugins(dates []string, plugins []cache.DataAdapter) {
	moduleName := "修复数据"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	barIndex++
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpRepair)
		bar.Add(1)
	}
	bar.Wait()
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}
