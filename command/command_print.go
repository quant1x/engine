package command

import (
	"fmt"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	cmder "github.com/spf13/cobra"
)

var (
	printModules = []cmdFlag[string]{}
)

const (
	printCommand     = "print"
	printDescription = "打印数据概要"
)

var (
	CmdPrint *cmder.Command = nil // CmdPrint 打印命令
)

func initPrint() {
	CmdPrint = &cmder.Command{
		Use:     printCommand,
		Example: Application + " print sh000001",
		Short:   printDescription,
		//Args:  cmder.MinimumNArgs(1),
		Run: func(cmd *cmder.Command, args []string) {
			tradeDate := cache.DefaultCanReadDate()
			if len(flagDate.Value) > 0 {
				tradeDate = exchange.FixTradeDate(flagDate.Value)
			}
			keywords := ""
			code := ""
			for _, m := range printModules {
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

	commandInit(CmdPrint, &flagDate)
	plugins := cache.Plugins(cache.PluginMaskFeature)
	printModules = make([]cmdFlag[string], len(plugins))
	for i, plugin := range plugins {
		key := plugin.Key()
		name := plugin.Name()
		printModules[i] = cmdFlag[string]{Name: key, Usage: plugin.Owner() + ": " + name, Value: ""}
		CmdPrint.Flags().StringVar(&(printModules[i].Value), printModules[i].Name, "", printModules[i].Usage)
	}
}

// 输出结构化信息
func handlePrintData(code, date string, plugin cache.DataAdapter) {
	securityCode := exchange.CorrectSecurityCode(code)
	name := securities.GetStockName(securityCode)
	tradeDate := exchange.FixTradeDate(date)
	cacheDate, featureDate := cache.CorrectDate(tradeDate)
	fmt.Printf("%s: %s, cache: %s, feature: %s\n", securityCode, name, cacheDate, featureDate)
	plugin.Print(securityCode, tradeDate)
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
	securityCode = exchange.CorrectSecurityCode(securityCode)
	name := securities.GetStockName(securityCode)
	fmt.Printf("%s: %s, %s\n", securityCode, name, tradeDate)
	df := factors.BasicKLine(securityCode)
	fmt.Println(df)
}
