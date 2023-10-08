package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/features"
	"gitee.com/quant1x/engine/flash"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/util/treemap"
)

// Repair 回补数据
func Repair(cacheDate, featureDate string) {
	allCodes := market.GetCodeList()
	for _, cache := range flash.CacheList() {
		//cache.FromHistory()
		fmt.Println(cache.Name())
		//ca := cache.Cache().(cachel5.CacheAdapter)
		mapFeature := treemap.NewWithStringComparator()
		codes := allCodes[:]
		for _, code := range codes {
			data := cache.Factory(cacheDate, code).(features.Feature)
			if data.Kind() != features.FeatureHistory {
				history := flash.GetL5History(code, cacheDate)
				if history != nil {
					data = data.FromHistory(*history)
				}
			}
			data.Repair(code, cacheDate, featureDate, true)
			mapFeature.Put(code, data)
		}
		// 加载缓存
		cache.Checkout(cacheDate)
		// 合并
		cache.Merge(mapFeature)
	}
}

// RepairAllFeature 回补更新特征
func RepairAllFeature(barIndex *int, cacheDate, featureDate string) {
	moduleName := "回补特征数据" + cacheDate
	allCodes := market.GetCodeList()
	cacheList := flash.CacheList()
	cacheCount := len(cacheList)
	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)
	for _, cache := range cacheList {
		mapFeature := treemap.NewWithStringComparator()
		codeCount := len(allCodes)
		barCode := progressbar.NewBar(*barIndex+1, "执行["+cache.Name()+"]", codeCount)
		dataSource := cache.Factory(featureDate, "")
		_ = dataSource.Init()
		for _, code := range allCodes {
			data := cache.Factory(cacheDate, code).(features.Feature)
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
		cache.Checkout(cacheDate)
		// 合并
		cache.Merge(mapFeature)
		_ = cache
		barCache.Add(1)
	}
	_ = allCodes
}
