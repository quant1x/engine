package models

import (
	"fmt"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/gox/util/treemap"
	"gitee.com/quant1x/pkg/tablewriter"
	"os"
)

// ExecuteStrategy 执行策略
func ExecuteStrategy(model Strategy, barIndex *int) {
	// 加载即时行情
	SyncAllSnapshots(barIndex)
	fmt.Println()
	// 执行策略
	allCodes := market.GetCodeList()
	count := len(allCodes)
	bar := progressbar.NewBar(*barIndex, "执行["+model.Name()+"]", count)
	results := treemap.NewWithStringComparator()
	for _, securityCode := range allCodes {
		// 此处可以增加过滤规则
		model.Evaluate(securityCode, results)
		bar.Add(1)
	}
	// 输出一个换行符, 结束上一个进度条
	fmt.Println()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tags.GetHeadersByTags(ResultInfo{}))
	results.Each(func(key interface{}, value interface{}) {
		row := value.(ResultInfo)
		table.Append(tags.GetValuesByTags(row))
	})
	table.Render()
}
