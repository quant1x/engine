package command

import (
	"strings"
	"unsafe"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/pkg/tablewriter"
	cmder "github.com/spf13/cobra"
)

const (
	defaultFlagAll = "all" // 全部
)

var (
	flagAll       = cmdFlag[bool]{Name: "all", Value: false, Usage: "全部"}
	flagBaseData  = cmdFlag[string]{Name: "base", Value: "", Usage: "基础数据"}
	flagFeatures  = cmdFlag[string]{Name: "features", Value: "", Usage: "特征数据"}
	flagStartDate = cmdFlag[string]{Name: "start", Value: exchange.LastTradeDate(), Usage: "开始日期"}
	flagEndDate   = cmdFlag[string]{Name: "end", Value: exchange.LastTradeDate(), Usage: "结束日期"}
	flagDate      = cmdFlag[string]{Name: "date", Value: "", Usage: "日期"}
)

type Command = cmder.Command

type cmdFlag[T ~int | ~bool | ~string] struct {
	Name  string
	Usage string
	Value T
}

func (cf *cmdFlag[T]) init(cmd *cmder.Command) {
	switch v := any(cf.Value).(type) {
	case bool:
		cmd.Flags().BoolVar((*bool)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case int:
		cmd.Flags().IntVar((*int)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case string:
		cmd.Flags().StringVar((*string)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	}
}

func commandInit[T ~int | ~bool | ~string](cmd *cmder.Command, cf *cmdFlag[T]) {
	switch v := any(cf.Value).(type) {
	case bool:
		cmd.Flags().BoolVar((*bool)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case int:
		cmd.Flags().IntVar((*int)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	case string:
		cmd.Flags().StringVar((*string)(unsafe.Pointer(&cf.Value)), cf.Name, v, cf.Usage)
	}
}

func parseFields(text string) (all bool, keywords []string) {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return
	}
	arr := strings.Split(text, ",")
	var list []string
	for _, v := range arr {
		v := strings.TrimSpace(v)
		if v == defaultFlagAll {
			return true, keywords
		} else if len(v) == 0 {
			continue
		} else {
			list = append(list, v)
		}
	}
	return false, list
}

type optionUsage struct {
	Key      string `name:"关键字"`
	Name     string `name:"概要"`
	Provider string `name:"数据源"`
	Usage    string `name:"描述"`
}

// 获取插件的Usage信息
func getPluginsUsage(plugins []cache.DataAdapter) string {
	writer := strings.Builder{}
	table := tablewriter.NewWriter(&writer)
	table.SetHeader(tags.GetHeadersByTags(optionUsage{}))
	table.Append(tags.GetValuesByTags(optionUsage{Key: defaultFlagAll, Name: "全部", Provider: ""}))
	for _, plugin := range plugins {
		ou := optionUsage{Key: plugin.Key(), Name: plugin.Name(), Provider: plugin.Owner(), Usage: plugin.Usage()}
		table.Append(tags.GetValuesByTags(ou))
	}
	table.EnableBorder(false)
	table.Render()
	return writer.String()
}
