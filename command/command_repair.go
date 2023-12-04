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

var (
	repairBases    = []cmdFlag[bool]{}
	repairFeatures = []cmdFlag[bool]{}
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
			for _, m := range repairBases {
				if m.Value {
					keywords = append(keywords, m.Name)
					break
				}
			}
			if len(keywords) == 0 {
				handleRepairAllDataSets(dates)
			} else {
				plugins := cache.PluginsWithName(cache.PluginMaskBaseData, keywords...)
				if len(plugins) == 0 {
					fmt.Printf("没有找到名字是[%s]的数据插件\n", strings.Join(keywords, ","))
				} else {
					handleRepairDataSetsWithPlugins(dates, plugins)
				}
			}
		} else if flagFeatures.Value {
			keywords := []string{}
			for _, m := range repairFeatures {
				if m.Value {
					keywords = append(keywords, m.Name)
					break
				}
			}
			if len(keywords) == 0 {
				handleRepairAllFeatures(dates)
			} else {
				plugins := cache.PluginsWithName(cache.PluginMaskFeature, keywords...)
				if len(plugins) == 0 {
					fmt.Printf("没有找到名字是[%s]的数据插件\n", strings.Join(keywords, ","))
				} else {
					handleRepairFeaturesWithPlugins(dates, plugins)
				}
			}

		}
	},
}

func initRepair() {
	commandInit(CmdRepair, &flagAll)
	commandInit(CmdRepair, &flagBaseData)
	commandInit(CmdRepair, &flagFeatures)
	commandInit(CmdRepair, &flagStartDate)
	commandInit(CmdRepair, &flagEndDate)

	plugins := cache.Plugins(cache.PluginMaskBaseData)
	repairBases = make([]cmdFlag[bool], len(plugins))
	for i, plugin := range plugins {
		key := plugin.Key()
		usage := plugin.Usage()
		repairBases[i] = cmdFlag[bool]{Name: key, Usage: plugin.Owner() + ": " + usage, Value: false}
		CmdRepair.Flags().BoolVar(&(repairBases[i].Value), repairBases[i].Name, repairBases[i].Value, repairBases[i].Usage)
	}

	plugins = cache.Plugins(cache.PluginMaskFeature)
	repairFeatures = make([]cmdFlag[bool], len(plugins))
	for i, plugin := range plugins {
		key := plugin.Key()
		usage := plugin.Usage()
		repairFeatures[i] = cmdFlag[bool]{Name: key, Usage: plugin.Owner() + ": " + usage, Value: false}
		CmdRepair.Flags().BoolVar(&(repairFeatures[i].Value), repairFeatures[i].Name, repairFeatures[i].Value, repairFeatures[i].Usage)
	}
}

func handleRepairAll(dates []string) {
	handleRepairAllDataSets(dates)
	handleRepairAllFeatures(dates)
}

func handleRepairAllDataSets(dates []string) {
	fmt.Println()
	moduleName := "补登数据集合"
	logger.Info(moduleName + ", 任务开始")
	mask := cache.PluginMaskBaseData
	plugins := cache.Plugins(mask)
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		//cacheDate, featureDate := cache.CorrectDate(date)
		barIndex++
		storages.BaseDataUpdate(barIndex, date, plugins, cache.OpRepair)
		bar.Add(1)
	}
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}

// 修复 - 指定的基础数据
func handleRepairDataSetsWithPlugins(dates []string, plugins []cache.DataAdapter) {
	fmt.Println()
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
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		barIndex++
		storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpRepair)
		bar.Add(1)
	}
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}

// 修复 - 指定的特征数据
func handleRepairFeaturesWithPlugins(dates []string, plugins []cache.DataAdapter) {
	fmt.Println()
	moduleName := "修复数据"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		barIndex++
		storages.FeaturesUpdate(&barIndex, cacheDate, featureDate, plugins, cache.OpRepair)
		bar.Add(1)
	}
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
}
