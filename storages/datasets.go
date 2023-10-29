package storages

import (
	"context"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/util/runtime"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/text/runewidth"
	"strings"
	"sync"
)

// 更新单个数据集
func updateOneDataSet(wg *sync.WaitGroup, parent, bar *progressbar.Bar, dataSet datasets.DataSet, date string, op cache.OpKind) {
	defer runtime.CatchPanic()
	moduleName := "基础数据"
	if op == cache.OpRepair {
		moduleName = "修复" + moduleName
	} else {
		moduleName = "更新" + moduleName
	}
	logger.Infof("%s: %s, begin", moduleName, dataSet.Name())
	allCodes := market.GetCodeList()
	for _, code := range allCodes {
		data := dataSet.Clone(date, code).(datasets.DataSet)
		if op == cache.OpUpdate {
			data.Update(date)
		} else if op == cache.OpRepair {
			data.Repair(date)
		}
		bar.Add(1)
	}
	parent.Add(1)
	wg.Done()
	logger.Infof("%s: %s, end", moduleName, dataSet.Name())
}

// BaseDataUpdate 修复数据
func BaseDataUpdate(barIndex int, date string, plugins []cache.DataAdapter, op cache.OpKind) {
	moduleName := "基础数据"
	if op == cache.OpRepair {
		moduleName = "修复" + moduleName
	} else {
		moduleName = "更新" + moduleName
	}
	var dataSetList []datasets.DataSet
	// 1.1 缓存数据集名称的最大宽度
	maxWidth := 0
	for _, plugin := range plugins {
		dataSet, ok := plugin.(datasets.DataSet)
		if ok {
			dataSetList = append(dataSetList, dataSet)
			width := runewidth.StringWidth(dataSet.Name())
			if width > maxWidth {
				maxWidth = width
			}
		}
	}
	logger.Infof("%s: all, begin", moduleName)
	// 2. 遍历全部数据插件
	dataSetCount := len(dataSetList)
	barCache := progressbar.NewBar(barIndex, "执行["+date+":"+moduleName+"]", dataSetCount)

	allCodes := market.GetCodeList()
	codeCount := len(allCodes)
	var wg sync.WaitGroup

	parent := coroutine.Context()
	ctx := context.WithValue(parent, cache.KBarIndex, barIndex)
	for sequence, dataSet := range dataSetList {
		_ = dataSet.Init(ctx, date)
		//format := fmt.Sprintf("%%%ds", maxWidth)
		//title := fmt.Sprintf(format, dataSet.Name())

		desc := dataSet.Name()
		width := runewidth.StringWidth(desc)
		title := strings.Repeat(" ", maxWidth-width) + desc
		barNo := barIndex + 1
		if cache.UseGoroutine {
			barNo += sequence
		}
		barCode := progressbar.NewBar(barNo, "执行["+title+"]", codeCount)
		wg.Add(1)
		if cache.UseGoroutine {
			go updateOneDataSet(&wg, barCache, barCode, dataSet, date, op)
		} else {
			updateOneDataSet(&wg, barCache, barCode, dataSet, date, op)
		}
	}
	wg.Wait()
	logger.Infof("%s: all, end", moduleName)
}
