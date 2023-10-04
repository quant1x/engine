package storages

import (
	"gitee.com/quant1x/engine/features"
	"gitee.com/quant1x/engine/flash"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/util/treemap"
)

// UpdateFeature 更新特征
func UpdateFeature(barIndex *int, cacheDate, featureDate string) {
	moduleName := "更新特征数据"
	allCodes := market.GetCodeList()
	cacheList := flash.CacheList()
	cacheCount := len(cacheList)
	//fmt.Println("\n\n")
	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)
	//barCache.Add(0)
	//fmt.Println()
	for _, cache := range cacheList {
		mapFeature := treemap.NewWithStringComparator()
		codeCount := len(allCodes)
		barCode := progressbar.NewBar(*barIndex+1, "执行["+cache.Name()+"]", codeCount)
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
