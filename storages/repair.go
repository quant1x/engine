package storages

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/datasets"
	"gitee.com/quant1x/engine/features"
	"gitee.com/quant1x/engine/flash"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/text/runewidth"
	"gitee.com/quant1x/gox/util/treemap"
	"strings"
	"sync"
)

// RepairBaseData 修复基础数据
func RepairBaseData(barIndex *int, cacheDate, featureDate string) {
	moduleName := "修复基础数据"
	// 1. 获取全部注册的数据集插件
	mask := cache.PluginMaskDataSet
	//dataSetList := flash.DataSetList()
	plugins := cache.Plugins(mask)
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

	// 2. 遍历全部数据插件
	dataSetCount := len(dataSetList)
	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", dataSetCount)

	allCodes := market.GetCodeList()
	var wg sync.WaitGroup

	for sequence, dataSet := range dataSetList {
		_ = dataSet.Init(barIndex, featureDate)
		codeCount := len(allCodes)
		//format := fmt.Sprintf("%%%ds", maxWidth)
		//title := fmt.Sprintf(format, dataSet.Name())
		width := runewidth.StringWidth(dataSet.Name())
		title := strings.Repeat(" ", maxWidth-width) + dataSet.Name()
		barNo := *barIndex + 1
		if useGoroutine {
			barNo += sequence
		}
		barCode := progressbar.NewBar(barNo, "执行["+title+"]", codeCount)
		wg.Add(1)
		if useGoroutine {
			go repairDateSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate)
		} else {
			repairDateSet(&wg, barCache, barCode, dataSet, cacheDate, featureDate)
		}
	}
	wg.Wait()
}

// RepairFeatures 修复特征数据
func RepairFeatures(barIndex *int, cacheDate, featureDate string) {
	moduleName := "修复特征数据" + cacheDate
	// 1. 获取全部注册的数据集插件
	mask := cache.PluginMaskFeature
	//dataSetList := flash.DataSetList()
	plugins := cache.Plugins(mask)
	var cacheList []cachel5.CacheAdapter
	maxWidth := 0
	for _, plugin := range plugins {
		adapter, ok := plugin.(cachel5.CacheAdapter)
		if ok {
			cacheList = append(cacheList, adapter)
			width := runewidth.StringWidth(adapter.Name())
			if width > maxWidth {
				maxWidth = width
			}
		}
	}
	allCodes := market.GetCodeList()
	cacheCount := len(cacheList)
	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)
	for _, adapter := range cacheList {
		mapFeature := treemap.NewWithStringComparator()
		codeCount := len(allCodes)
		barCode := progressbar.NewBar(*barIndex+1, "执行["+adapter.Name()+"]", codeCount)
		dataSource := adapter.Factory(featureDate, "")
		_ = dataSource.Init(barIndex, featureDate)
		for _, code := range allCodes {
			data := adapter.Factory(cacheDate, code).(features.Feature)
			if data.Kind() != features.FeatureHistory {
				history := flash.GetL5History(code, cacheDate)
				if history != nil {
					data = data.FromHistory(*history)
				}
			}
			data.Repair(code, cacheDate, featureDate, true)
			mapFeature.Put(code, data)
			barCode.Add(1)
		}
		// 加载缓存
		adapter.Checkout(cacheDate)
		// 合并
		adapter.Merge(mapFeature)
		_ = adapter
		barCache.Add(1)
	}
	_ = allCodes
}
