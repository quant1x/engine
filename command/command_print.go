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
	tradeDate    = ""
	f10Code      = "" // F10
	exchangeCode = "" // Exchange
	maCode       = "" // 移动平均线
	boxCode      = "" // 平台
)

// CmdPrint 打印命令
var CmdPrint = &cmder.Command{
	Use:     "print",
	Example: Application + " print sh000001",
	Short:   "打印K线概要",
	//Args:  cmder.MinimumNArgs(1),
	Run: func(cmd *cmder.Command, args []string) {
		tradeDate = trading.FixTradeDate(tradeDate)
		if len(exchangeCode) > 0 {
			//printExchange(exchangeCode)
		} else if len(f10Code) > 0 {
			printF10(f10Code)
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
			printKline(securityCode)
		}
	},
}

func init() {
	CmdPrint.Flags().StringVar(&f10Code, "f10", "", "查看快照扩展数据")
	CmdPrint.Flags().StringVar(&exchangeCode, "exchange", "", "查看快照扩展数据")
	CmdPrint.Flags().StringVar(&maCode, "ma", "", "查看均线")
	CmdPrint.Flags().StringVar(&boxCode, "box", "", "查看平台数据")
	CmdPrint.Flags().StringVar(&tradeDate, "date", cache.DefaultCanReadDate(), "指定日期")
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
func printKline(securityCode string) {
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

func printF10(securityCode string) {
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
