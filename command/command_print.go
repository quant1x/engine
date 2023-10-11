package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gotdx/proto"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/gotdx/trading"
	"gitee.com/quant1x/gox/tags"
	"github.com/olekukonko/tablewriter"
	cmder "github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	exchangeCode = "" // Exchange
	maCode       = "" // 移动平均线
	boxCode      = "" // 平台
	subModules   = []cmdFlag[string]{}
)

// CmdPrint 打印命令
var CmdPrint = &cmder.Command{
	Use:     "print",
	Example: Application + " print sh000001",
	Short:   "打印K线概要",
	//Args:  cmder.MinimumNArgs(1),
	Run: func(cmd *cmder.Command, args []string) {
		tradeDate := cache.DefaultCanReadDate()
		if len(flagDate.Value) > 0 {
			tradeDate = trading.FixTradeDate(flagDate.Value)
		}
		keywords := ""
		code := ""
		for _, m := range subModules {
			if len(m.Value) > 0 {
				keywords = m.Name
				code = m.Value
				break
			}
		}
		if len(keywords) > 0 {
			plugins := cache.PluginsWithName(cache.PluginMaskFeature, keywords)
			if len(plugins) == 0 {
				fmt.Printf("没有找到名字是[%s]的数据插件\n", keywords)
			} else {
				handlePrintData(code, tradeDate, plugins[0])
			}
		} else if len(maCode) > 0 {
			//printMA(maCode)
		} else if len(boxCode) > 0 {
			//printBox(boxCode)
		} else {
			if len(args) != 1 {
				fmt.Println(cmd.Help())
				return
			}
			// 默认输出K线信息
			securityCode := args[0]
			printKline(securityCode, tradeDate)
		}
	},
}

func init() {
	commandInit(CmdPrint, &flagDate)
	plugins := cache.Plugins(cache.PluginMaskFeature)
	subModules = make([]cmdFlag[string], len(plugins))
	for i, plugin := range plugins {
		key := plugin.Key()
		usage := plugin.Usage()
		subModules[i] = cmdFlag[string]{Name: key, Usage: usage, Value: ""}
		CmdPrint.Flags().StringVar(&(subModules[i].Value), subModules[i].Name, "", subModules[i].Usage)
	}
	//CmdPrint.Flags().StringVar(&f10Code, "f10", "", "查看快照扩展数据")
	//CmdPrint.Flags().StringVar(&exchangeCode, "exchange", "", "查看快照扩展数据")
	//CmdPrint.Flags().StringVar(&maCode, "ma", "", "查看均线")
	//CmdPrint.Flags().StringVar(&boxCode, "box", "", "查看平台数据")
	//CmdPrint.Flags().StringVar(&tradeDate, "date", cache.DefaultCanReadDate(), "指定日期")
}

// 输出结构化信息
func handlePrintData(code, date string, plugin cache.DataPlugin) {
	fmt.Println()
	securityCode := proto.CorrectSecurityCode(code)
	name := securities.GetStockName(securityCode)
	tradeDate := trading.FixTradeDate(date)
	fmt.Printf("%s: %s, %s\n", securityCode, name, tradeDate)
	plugin.Print(securityCode, tradeDate)

	//fmt.Println()
}

//// 输出K线概要数据列表
//func v1PrintKline(securityCode string) {
//	securityCode = proto.CorrectSecurityCode(securityCode)
//	name := securities.GetStockName(securityCode)
//	fmt.Printf("%s: %s\n", securityCode, name)
//	df := datasets.KLine(securityCode)
//	fmt.Println(df)
//}

// 输出K线概要数据列表
func printKline(securityCode string, tradeDate string) {
	securityCode = proto.CorrectSecurityCode(securityCode)
	name := securities.GetStockName(securityCode)
	fmt.Printf("%s: %s, %s\n", securityCode, name, tradeDate)
	df := datasets.BasicKLine(securityCode)
	fmt.Println(df)
}

func checkoutTable(v any) (headers []string, records [][]string) {
	headers = []string{"字段", "数值"}
	fields := tags.GetHeadersByTags(v)
	values := tags.GetValuesByTags(v)
	num := len(fields)
	if num > len(values) {
		num = len(values)
	}
	for i := 0; i < num; i++ {
		records = append(records, []string{fields[i], strings.TrimSpace(values[i])})
	}
	return
}

func printF10(securityCode string, tradeDate string) {
	securityCode = proto.CorrectSecurityCode(securityCode)
	name := securities.GetStockName(securityCode)
	fmt.Printf("%s: %s, %s\n", securityCode, name, tradeDate)
	value := smart.GetL5F10(securityCode, tradeDate)
	headers, records := checkoutTable(value)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	table.AppendBulk(records)
	table.Render()
}

//func printExchange(securityCode string) {
//	securityCode = proto.CorrectSecurityCode(securityCode)
//	name := securities.GetStockName(securityCode)
//	fmt.Printf("%s: %s, %s\n", securityCode, name, tradeDate)
//	value := flash.GetL5Exchange(securityCode, tradeDate)
//	headers, records := checkoutTable(value)
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader(headers)
//	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
//	table.AppendBulk(records)
//	table.Render()
//}
//
//func printMA(securityCode string) {
//	securityCode = proto.CorrectSecurityCode(securityCode)
//	name := securities.GetStockName(securityCode)
//	fmt.Printf("%s: %s, %s\n", securityCode, name, tradeDate)
//	value := flash.GetL5MovingAverage(securityCode, tradeDate)
//	headers, records := checkoutTable(value)
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader(headers)
//	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
//	table.AppendBulk(records)
//	table.Render()
//}
//
//func printBox(securityCode string) {
//	securityCode = proto.CorrectSecurityCode(securityCode)
//	name := securities.GetStockName(securityCode)
//	fmt.Printf("%s: %s, %s\n", securityCode, name, tradeDate)
//	value := flash.GetL5Box(securityCode, tradeDate)
//	headers, records := checkoutTable(value)
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader(headers)
//	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
//	table.AppendBulk(records)
//	table.Render()
//}
