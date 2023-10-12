package storages

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/text/runewidth"
	"gitee.com/quant1x/gox/util/treemap"
)

// FeaturesUpdate 更新特征
func FeaturesUpdate(barIndex *int, cacheDate, featureDate string, plugins []cache.DataPlugin, op cache.OpKind) {
	moduleName := "特征数据"
	if op == cache.OpRepair {
		moduleName = "修复" + moduleName
	} else {
		moduleName = "更新" + moduleName
	}
	moduleName += cacheDate
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
			data := adapter.Factory(cacheDate, code).(factors.Feature)
			if data.Kind() != factors.FeatureHistory {
				history := smart.GetL5History(code, cacheDate)
				if history != nil {
					data = data.FromHistory(*history)
				}
			}
			if op == cache.OpRepair {
				data.Repair(code, cacheDate, featureDate, true)
			} else {
				data.Update(code, cacheDate, featureDate, true)
			}
			mapFeature.Put(code, data)
			barCode.Add(1)
		}
		// 加载缓存
		adapter.Checkout(cacheDate)
		// 合并
		adapter.Merge(mapFeature)
		barCache.Add(1)
	}
}

//// RepairFeatures 修复特征数据
//func RepairFeatures(barIndex *int, cacheDate, featureDate string) {
//	moduleName := "修复特征数据" + cacheDate
//	// 1. 获取全部注册的数据集插件
//	mask := cache.PluginMaskFeature
//	//dataSetList := flash.DataSetList()
//	plugins := cache.Plugins(mask)
//	var cacheList []cachel5.CacheAdapter
//	maxWidth := 0
//	for _, plugin := range plugins {
//		adapter, ok := plugin.(cachel5.CacheAdapter)
//		if ok {
//			cacheList = append(cacheList, adapter)
//			width := runewidth.StringWidth(adapter.Name())
//			if width > maxWidth {
//				maxWidth = width
//			}
//		}
//	}
//	allCodes := market.GetCodeList()
//	cacheCount := len(cacheList)
//	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)
//	for _, adapter := range cacheList {
//		mapFeature := treemap.NewWithStringComparator()
//		codeCount := len(allCodes)
//		barCode := progressbar.NewBar(*barIndex+1, "执行["+adapter.Name()+"]", codeCount)
//		dataSource := adapter.Factory(featureDate, "")
//		_ = dataSource.Init(barIndex, featureDate)
//		for _, code := range allCodes {
//			data := adapter.Factory(cacheDate, code).(factors.Feature)
//			if data.Kind() != factors.FeatureHistory {
//				history := smart.GetL5History(code, cacheDate)
//				if history != nil {
//					data = data.FromHistory(*history)
//				}
//			}
//			data.Repair(code, cacheDate, featureDate, true)
//			mapFeature.Put(code, data)
//			barCode.Add(1)
//		}
//		// 加载缓存
//		adapter.Checkout(cacheDate)
//		// 合并
//		adapter.Merge(mapFeature)
//		_ = adapter
//		barCache.Add(1)
//	}
//	_ = allCodes
//}
