package storages

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/datasets"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/text/runewidth"
	"strings"
	"sync"
)

//// BaseDataUpdate 修复基础数据
//func BaseDataUpdate(barIndex *int, cacheDate, featureDate string) {
//	moduleName := "修复基础数据"
//	// 1. 获取全部注册的数据集插件
//	mask := cache.PluginMaskBaseData
//	//dataSetList := flash.DataSetList()
//	plugins := cache.Plugins(mask)
//	var dataSetList []datasets.DataSet
//	// 1.1 缓存数据集名称的最大宽度
//	maxWidth := 0
//	for _, plugin := range plugins {
//		dataSet, ok := plugin.(datasets.DataSet)
//		if ok {
//			dataSetList = append(dataSetList, dataSet)
//			width := runewidth.StringWidth(dataSet.Name())
//			if width > maxWidth {
//				maxWidth = width
//			}
//		}
//	}
//
//	// 2. 遍历全部数据插件
//	dataSetCount := len(dataSetList)
//	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", dataSetCount)
//
//	allCodes := market.GetCodeList()
//	var wg sync.WaitGroup
//
//	for sequence, dataSet := range dataSetList {
//		_ = dataSet.Init(barIndex, featureDate)
//		codeCount := len(allCodes)
//		//format := fmt.Sprintf("%%%ds", maxWidth)
//		//title := fmt.Sprintf(format, dataSet.Name())
//		width := runewidth.StringWidth(dataSet.Name())
//		title := strings.Repeat(" ", maxWidth-width) + dataSet.Name()
//		barNo := *barIndex + 1
//		if useGoroutine {
//			barNo += sequence
//		}
//		barCode := progressbar.NewBar(barNo, "执行["+title+"]", codeCount)
//		wg.Add(1)
//		if useGoroutine {
//			go updateOneDataSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate, opRepair)
//		} else {
//			updateOneDataSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate, opRepair)
//		}
//	}
//	wg.Wait()
//}

// 更新单个数据集
func updateOneDataSet(wg *sync.WaitGroup, parent, bar *progressbar.Bar, dataSet datasets.DataSet, cacheDate, featureDate string, op cache.OpKind) {
	allCodes := market.GetCodeList()
	for _, code := range allCodes {
		data := dataSet.Clone(cacheDate, code).(datasets.DataSet)
		if op == cache.OpUpdate {
			data.Update(cacheDate, featureDate)
		} else if op == cache.OpRepair {
			data.Repair(cacheDate, featureDate)
		}
		bar.Add(1)
	}
	parent.Add(1)
	wg.Done()
}

// BaseDataUpdate 修复数据
func BaseDataUpdate(barIndex int, cacheDate, featureDate string, plugins []cache.DataPlugin, op cache.OpKind) {
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
			width := runewidth.StringWidth(dataSet.Desc())
			if width > maxWidth {
				maxWidth = width
			}
		}
	}

	// 2. 遍历全部数据插件
	dataSetCount := len(dataSetList)
	barCache := progressbar.NewBar(barIndex, "执行["+cacheDate+":"+moduleName+"]", dataSetCount)

	allCodes := market.GetCodeList()
	codeCount := len(allCodes)
	var wg sync.WaitGroup

	for sequence, dataSet := range dataSetList {
		_ = dataSet.Init(&barIndex, featureDate)
		//format := fmt.Sprintf("%%%ds", maxWidth)
		//title := fmt.Sprintf(format, dataSet.Name())
		desc := dataSet.Desc()
		width := runewidth.StringWidth(desc)
		title := strings.Repeat(" ", maxWidth-width) + desc
		barNo := barIndex + 1
		if cache.UseGoroutine {
			barNo += sequence
		}
		barCode := progressbar.NewBar(barNo, "执行["+title+"]", codeCount)
		wg.Add(1)
		if cache.UseGoroutine {
			go updateOneDataSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate, op)
		} else {
			updateOneDataSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate, op)
		}
	}
	wg.Wait()
}

//// UpdateBaseData 更新基础数据
//func UpdateBaseData(barIndex *int, cacheDate, featureDate string) {
//	moduleName := "更新基础数据"
//	// 1. 获取全部注册的数据集插件
//	mask := cache.PluginMaskBaseData
//	//dataSetList := flash.DataSetList()
//	plugins := cache.Plugins(mask)
//	var dataSetList []datasets.DataSet
//	// 1.1 缓存数据集名称的最大宽度
//	maxWidth := 0
//	for _, plugin := range plugins {
//		dataSet, ok := plugin.(datasets.DataSet)
//		if ok {
//			dataSetList = append(dataSetList, dataSet)
//			width := runewidth.StringWidth(dataSet.Name())
//			if width > maxWidth {
//				maxWidth = width
//			}
//		}
//	}
//
//	// 2. 遍历全部数据插件
//	dataSetCount := len(dataSetList)
//	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", dataSetCount)
//
//	allCodes := market.GetCodeList()
//	var wg sync.WaitGroup
//
//	for sequence, dataSet := range dataSetList {
//		_ = dataSet.Init(barIndex, featureDate)
//		codeCount := len(allCodes)
//		//format := fmt.Sprintf("%%%ds", maxWidth)
//		//title := fmt.Sprintf(format, dataSet.Name())
//		width := runewidth.StringWidth(dataSet.Name())
//		title := strings.Repeat(" ", maxWidth-width) + dataSet.Name()
//		barNo := *barIndex + 1
//		if useGoroutine {
//			barNo += sequence
//		}
//		barCode := progressbar.NewBar(barNo, "执行["+title+"]", codeCount)
//		wg.Add(1)
//		if useGoroutine {
//			go updateOneDataSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate, opUpdate)
//		} else {
//			updateOneDataSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate, opUpdate)
//		}
//	}
//	wg.Wait()
//}
