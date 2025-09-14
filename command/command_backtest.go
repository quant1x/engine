package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/quant1x/engine/cache"
	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/engine/storages"
	"github.com/quant1x/exchange"
	"github.com/quant1x/pkg/tablewriter"
	"github.com/quant1x/x/logger"
	"github.com/quant1x/x/progressbar"
	"github.com/quant1x/x/tags"
	cli "github.com/spf13/cobra"
)

const (
	backTestCommand     = "backtest"
	backTestDescription = "回测"
)

var (
	// 回测
	cmdBackTest *cli.Command = nil
)

func initBackTest() {
	cmdBackTest = &cli.Command{
		Use:     backTestCommand,
		Example: Application + " " + backTestCommand + " --all",
		//Args:    args.MinimumNArgs(0),
		Args: func(cmd *cli.Command, args []string) error {
			return nil
		},
		Short: backTestDescription,
		Long:  backTestDescription,
		Run: func(cmd *cli.Command, args []string) {
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
			// 连续2个空白行
			consecutiveEmptyLines := strings.Repeat("\r\n", 2)
			fmt.Printf("%s数据: %s => %s"+consecutiveEmptyLines, backTestDescription, dates[0], dates[count-1])
			base.UpdateBeginDateOfHistoricalTradingData(dates[0])
			if flagAll.Value {
				//handleBacktestAll(dates)
			} else if len(flagBaseData.Value) > 0 {
				//all, keywords := parseFields(flagBaseData.Value)
				//if all || len(keywords) == 0 {
				//	handleBacktestAllDataSets(dates)
				//} else {
				//	plugins := cache.PluginsWithName(cache.PluginMaskBaseData, keywords...)
				//	if len(plugins) == 0 {
				//		fmt.Printf("没有找到名字是[%s]的数据插件\n", strings.Join(keywords, ","))
				//	} else {
				//		handleBacktestDataSetsWithPlugins(dates, plugins)
				//	}
				//}
			} else if len(flagFeatures.Value) > 0 {
				all, keywords := parseFields(flagFeatures.Value)
				if all || len(keywords) == 0 {
					//handleBacktestAllFeatures(dates)
				} else {
					plugins := cache.PluginsWithName(cache.PluginMaskFeature, keywords...)
					if len(plugins) == 0 {
						fmt.Printf("没有找到名字是[%s]的数据插件\n", strings.Join(keywords, ","))
					} else {
						handleBacktestFeaturesWithPlugins(dates, plugins)
					}
				}
			} else {
				fmt.Printf("Error: 非全部%s, 必须携带--features或--base\n", backTestDescription)
				_ = cmd.Usage()
			}

		},
	}
	commandInit(cmdBackTest, &flagAll)
	commandInit(cmdBackTest, &flagStartDate)
	commandInit(cmdBackTest, &flagEndDate)

	// 1. 基础数据
	plugins := cache.Plugins(cache.PluginMaskBaseData)
	flagBaseData.Usage = getPluginsUsage(plugins)
	commandInit(cmdBackTest, &flagBaseData)

	// 2. 特征数据
	plugins = cache.Plugins(cache.PluginMaskFeature)
	flagFeatures.Usage = getPluginsUsage(plugins)
	commandInit(cmdBackTest, &flagFeatures)
}

// 回测 - 指定的特征数据
func handleBacktestFeaturesWithPlugins(dates []string, plugins []cache.DataAdapter) {
	moduleName := "回测数据"
	logger.Info(moduleName + ", 任务开始")
	count := len(dates)
	barIndex := 1
	bar := progressbar.NewBar(barIndex, "执行["+moduleName+"]", count)
	barIndex++
	//var metrics []cache.AdapterMetric
	var rows [][]string
	for _, date := range dates {
		cacheDate, featureDate := cache.CorrectDate(date)
		result := storages.FeaturesBackTest(&barIndex, cacheDate, featureDate, plugins, cache.OpBackTest)
		for _, v := range result {
			row := []string{date}
			cols := tags.GetValuesByTags(v)
			row = append(row, cols...)
			rows = append(rows, row)
		}
		//if len(result) > 0 {
		//	metrics = append(metrics, result...)
		//}

		bar.Add(1)
	}
	bar.Wait()
	logger.Info(moduleName+", 任务执行完毕.", time.Now())
	fmt.Println()
	//metricCount := len(metrics)
	//if metricCount > 0 {
	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"date"}
	headers = append(headers, tags.GetHeadersByTags(cache.AdapterMetric{})...)
	table.SetAutoFormatHeaders(false)
	//table.SetAutoMergeCells(true)
	table.SetHeader(headers)
	//for i, v := range metrics {
	//	table.Append(tags.GetValuesByTags(v))
	//}
	table.AppendBulk(rows)
	table.Render()
	//}
}
