package storages

import (
	"fmt"
	"gitee.com/quant1x/engine/features"
	"gitee.com/quant1x/engine/flash"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/gox/util/treemap"
)

func Repair(cacheDate, featureDate string) {
	allCodes := market.GetCodeList()
	for _, cache := range flash.CacheList() {
		//cache.FromHistory()
		fmt.Println(cache.Name())
		//ca := cache.Cache().(cachel5.CacheAdapter)
		mapFeature := treemap.NewWithStringComparator()
		for _, code := range allCodes {
			data := cache.Factory(cacheDate, code).(features.Feature)
			if data.Kind() != features.FeatureHistory {
				history := flash.GetL5History(code, cacheDate)
				if history != nil {
					data = data.FromHistory(*history)
				}
			}
			data.Repair(cacheDate, featureDate)
			mapFeature.Put(code, data)
		}
		// 加载缓存
		cache.Checkout(cacheDate)
		// 合并
		cache.Merge(mapFeature)
	}
}
