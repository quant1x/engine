package storages

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/text/runewidth"
	"gitee.com/quant1x/gox/util/treemap"
	"github.com/panjf2000/ants/v2"
	"sync"
)

//// 更新单个特征
//func updateV1OneFeature(wg *sync.WaitGroup, parent, bar *progressbar.Bar, adapter cachel5.CacheAdapter, cacheDate, featureDate string, op cache.OpKind) {
//	allCodes := market.GetCodeList()
//	for _, code := range allCodes {
//		feature := adapter.Clone(cacheDate, code).(datasets.DataSet)
//		if op == cache.OpUpdate {
//			feature.Update(cacheDate, featureDate)
//		} else if op == cache.OpRepair {
//			feature.Repair(cacheDate, featureDate)
//		}
//		bar.Add(1)
//	}
//	parent.Add(1)
//	wg.Done()
//}

func updateStockFeature(wg *sync.WaitGroup, bar *progressbar.Bar, feature factors.Feature, code string, cacheDate, featureDate string, op cache.OpKind, p *treemap.Map) {
	if op == cache.OpRepair {
		feature.Repair(code, cacheDate, featureDate, true)
	} else {
		feature.Update(code, cacheDate, featureDate, true)
	}
	p.Put(code, feature)
	bar.Add(1)
	wg.Done()
}

// 更新单个特征
func updateOneFeature(parent, bar *progressbar.Bar, adapter cachel5.CacheAdapter, cacheDate, featureDate string, op cache.OpKind, barIndex *int) {
	mapFeature := treemap.NewWithStringComparator()
	var wg sync.WaitGroup
	dataSource := adapter.Factory(featureDate, "")
	_ = dataSource.Init(barIndex, featureDate)
	allCodes := market.GetCodeList()
	for _, code := range allCodes {
		feature := adapter.Factory(cacheDate, code).(factors.Feature)
		if feature.Kind() != factors.FeatureHistory {
			history := smart.GetL5History(code, cacheDate)
			if history != nil {
				feature = feature.FromHistory(*history)
			}
		}
		wg.Add(1)
		//if op == cache.OpRepair {
		//	feature.Repair(code, cacheDate, featureDate, true)
		//} else {
		//	feature.Update(code, cacheDate, featureDate, true)
		//}
		//mapFeature.Put(code, feature)
		//bar.Add(1)
		//if cache.UseGoroutine {
		//
		//}
		go updateStockFeature(&wg, bar, feature, code, cacheDate, featureDate, op, mapFeature)
	}
	// 加载缓存
	adapter.Checkout(cacheDate)
	// 合并
	adapter.Merge(mapFeature)
	wg.Wait()
	parent.Add(1)
}

//const (
//	queueMax = 4096
//)
//
//type featureTask struct {
//	wg           *sync.WaitGroup
//	bar          *progressbar.Bar
//	feature         factors.Feature
//	securityCode string
//	cacheDate    string
//	featureDate  string
//}
//
//// 更新单个特征
//func v2updateOneFeature(parent, bar *progressbar.Bar, adapter cachel5.CacheAdapter, cacheDate, featureDate string, op cache.OpKind, barIndex *int) {
//	mapFeature := treemap.NewWithStringComparator()
//	var wg sync.WaitGroup
//	dataSource := adapter.Factory(featureDate, "")
//	_ = dataSource.Init(barIndex, featureDate)
//	allCodes := market.GetCodeList()
//	queue := fastqueue.NewQueue[featureTask](queueMax)
//	queue.SetReadEvent(func(feature []featureTask) {
//		for _, v := range feature {
//			//v.wg.Add(1)
//			updateStockFeature(v.wg, v.bar, v.feature, v.securityCode, v.cacheDate, v.featureDate, op, mapFeature)
//			//fmt.Println(v)
//			//bar.Add(1)
//		}
//	})
//	for _, code := range allCodes {
//		feature := adapter.Factory(cacheDate, code).(factors.Feature)
//		if feature.Kind() != factors.FeatureHistory {
//			history := smart.GetL5History(code, cacheDate)
//			if history != nil {
//				feature = feature.FromHistory(*history)
//			}
//		}
//		wg.Add(1)
//		go queue.Push(featureTask{
//			wg:           &wg,
//			bar:          bar,
//			feature:         feature,
//			securityCode: code,
//			cacheDate:    cacheDate,
//			featureDate:  featureDate,
//		})
//	}
//	queue.Wait()
//	// 加载缓存
//	adapter.Checkout(cacheDate)
//	// 合并
//	adapter.Merge(mapFeature)
//	//wg.Wait()
//	parent.Add(1)
//}

type featureTask struct {
	wg           *sync.WaitGroup
	bar          *progressbar.Bar
	feature      factors.Feature
	securityCode string
	cacheDate    string
	featureDate  string
}

// FeaturesUpdate 更新特征
func FeaturesUpdate(barIndex *int, cacheDate, featureDate string, plugins []cache.DataPlugin, op cache.OpKind) {
	moduleName := "特征数据"
	if op == cache.OpRepair {
		moduleName = "修复" + moduleName
	} else {
		moduleName = "更新" + moduleName
	}
	moduleName += cacheDate
	var adapters []cachel5.CacheAdapter
	maxWidth := 0
	for _, plugin := range plugins {
		adapter, ok := plugin.(cachel5.CacheAdapter)
		if ok {
			adapters = append(adapters, adapter)
			width := runewidth.StringWidth(adapter.Name())
			if width > maxWidth {
				maxWidth = width
			}
		}
	}
	var wgAdapter sync.WaitGroup
	cacheCount := len(adapters)
	barAdapter := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)

	allCodes := market.GetCodeList()
	allCodes = allCodes[:]
	codeCount := len(allCodes)
	for _, adapter := range adapters {
		wgAdapter.Add(1)

		barCode := progressbar.NewBar(*barIndex+1, "执行["+adapter.Name()+"]", codeCount)
		//updateOneFeature(barAdapter, barCode, adapter, cacheDate, featureDate, op, barIndex)

		mapFeature := treemap.NewWithStringComparator()
		var wg sync.WaitGroup
		p, _ := ants.NewPoolWithFunc(quotes.POOL_MAX, func(i interface{}) {
			v := i.(featureTask)
			updateStockFeature(v.wg, v.bar, v.feature, v.securityCode, v.cacheDate, v.featureDate, op, mapFeature)
		})
		defer p.Release()

		dataSource := adapter.Factory(featureDate, "")
		_ = dataSource.Init(barIndex, featureDate)
		for _, code := range allCodes {
			feature := adapter.Factory(cacheDate, code).(factors.Feature)
			if feature.Kind() != factors.FeatureHistory {
				history := smart.GetL5History(code, cacheDate)
				if history != nil {
					feature = feature.FromHistory(*history)
				}
			}
			wg.Add(1)
			//go updateStockFeature(&wg, barCode, feature, code, cacheDate, featureDate, op, mapFeature)
			_ = p.Invoke(featureTask{
				wg:           &wg,
				bar:          barCode,
				feature:      feature,
				securityCode: code,
				cacheDate:    cacheDate,
				featureDate:  featureDate,
			})
		}
		wg.Wait()
		// 加载缓存
		adapter.Checkout(cacheDate)
		// 合并
		adapter.Merge(mapFeature)
		// 适配器进度条+1
		barAdapter.Add(1)
		wgAdapter.Done()
	}
	wgAdapter.Wait()
}
