package storages

import (
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gox/fastqueue"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/text/runewidth"
	"gitee.com/quant1x/gox/util/treemap"
	"sync"
)

//// 更新单个特征
//func updateV1OneFeature(wg *sync.WaitGroup, parent, bar *progressbar.Bar, adapter cachel5.CacheAdapter, cacheDate, featureDate string, op cache.OpKind) {
//	allCodes := market.GetCodeList()
//	for _, code := range allCodes {
//		data := adapter.Clone(cacheDate, code).(datasets.DataSet)
//		if op == cache.OpUpdate {
//			data.Update(cacheDate, featureDate)
//		} else if op == cache.OpRepair {
//			data.Repair(cacheDate, featureDate)
//		}
//		bar.Add(1)
//	}
//	parent.Add(1)
//	wg.Done()
//}

func updateStockFeature(wg *sync.WaitGroup, bar *progressbar.Bar, data factors.Feature, code string, cacheDate, featureDate string, op cache.OpKind, p *treemap.Map) {
	if op == cache.OpRepair {
		data.Repair(code, cacheDate, featureDate, true)
	} else {
		data.Update(code, cacheDate, featureDate, true)
	}
	p.Put(code, data)
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
		data := adapter.Factory(cacheDate, code).(factors.Feature)
		if data.Kind() != factors.FeatureHistory {
			history := smart.GetL5History(code, cacheDate)
			if history != nil {
				data = data.FromHistory(*history)
			}
		}
		wg.Add(1)
		//if op == cache.OpRepair {
		//	data.Repair(code, cacheDate, featureDate, true)
		//} else {
		//	data.Update(code, cacheDate, featureDate, true)
		//}
		//mapFeature.Put(code, data)
		//bar.Add(1)
		//if cache.UseGoroutine {
		//
		//}
		updateStockFeature(&wg, bar, data, code, cacheDate, featureDate, op, mapFeature)
	}
	// 加载缓存
	adapter.Checkout(cacheDate)
	// 合并
	adapter.Merge(mapFeature)
	wg.Wait()
	parent.Add(1)
}

const (
	queueMax = 4096
)

type featureTask struct {
	wg           *sync.WaitGroup
	bar          *progressbar.Bar
	data         factors.Feature
	securityCode string
	cacheDate    string
	featureDate  string
}

// 更新单个特征
func v2updateOneFeature(parent, bar *progressbar.Bar, adapter cachel5.CacheAdapter, cacheDate, featureDate string, op cache.OpKind, barIndex *int) {
	mapFeature := treemap.NewWithStringComparator()
	var wg sync.WaitGroup
	dataSource := adapter.Factory(featureDate, "")
	_ = dataSource.Init(barIndex, featureDate)
	allCodes := market.GetCodeList()
	queue := fastqueue.NewQueue[featureTask](queueMax)
	queue.SetReadEvent(func(data []featureTask) {
		for _, v := range data {
			//v.wg.Add(1)
			updateStockFeature(v.wg, v.bar, v.data, v.securityCode, v.cacheDate, v.featureDate, op, mapFeature)
			//fmt.Println(v)
			//bar.Add(1)
		}
	})
	for _, code := range allCodes {
		data := adapter.Factory(cacheDate, code).(factors.Feature)
		if data.Kind() != factors.FeatureHistory {
			history := smart.GetL5History(code, cacheDate)
			if history != nil {
				data = data.FromHistory(*history)
			}
		}
		wg.Add(1)
		go queue.Push(featureTask{
			wg:           &wg,
			bar:          bar,
			data:         data,
			securityCode: code,
			cacheDate:    cacheDate,
			featureDate:  featureDate,
		})
	}
	queue.Wait()
	// 加载缓存
	adapter.Checkout(cacheDate)
	// 合并
	adapter.Merge(mapFeature)
	//wg.Wait()
	parent.Add(1)
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
	var topWaitGroup sync.WaitGroup
	allCodes := market.GetCodeList()
	cacheCount := len(cacheList)
	barCache := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)
	for _, adapter := range cacheList {
		topWaitGroup.Add(1)
		codeCount := len(allCodes)
		barCode := progressbar.NewBar(*barIndex+1, "执行["+adapter.Name()+"]", codeCount)
		updateOneFeature(barCache, barCode, adapter, cacheDate, featureDate, op, barIndex)
		topWaitGroup.Done()
	}
	topWaitGroup.Wait()
}
