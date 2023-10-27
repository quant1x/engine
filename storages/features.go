package storages

import (
	"context"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/text/runewidth"
	"gitee.com/quant1x/gox/util/treemap"
	"runtime/debug"
	"sync"
	"time"
)

func updateStockFeature(wg *sync.WaitGroup, bar *progressbar.Bar, feature factors.Feature, code string, cacheDate, featureDate string, op cache.OpKind, p *treemap.Map, sb *cache.ScoreBoard) {
	defer func() {
		if err := recover(); err != nil {
			s := string(debug.Stack())
			logger.Errorf("err=%v, stack=%s", err, s)
		}
	}()
	now := time.Now()
	defer sb.Add(1, time.Since(now))
	if op == cache.OpRepair {
		feature.Repair(code, cacheDate, featureDate, true)
	} else {
		feature.Update(code, cacheDate, featureDate, true)
	}
	p.Put(code, feature)
	bar.Add(1)
	wg.Done()
}

// FeaturesUpdate 更新特征
func FeaturesUpdate(barIndex *int, cacheDate, featureDate string, plugins []cache.DataAdapter, op cache.OpKind) {
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
	logger.Infof("%s: all, begin", moduleName)

	var wgAdapter sync.WaitGroup
	cacheCount := len(adapters)
	barAdapter := progressbar.NewBar(*barIndex, "执行["+moduleName+"]", cacheCount)
	allCodes := market.GetCodeList()
	allCodes = allCodes[:]
	codeCount := len(allCodes)
	for _, adapter := range adapters {
		logger.Infof("%s: %s, begin", moduleName, adapter.Name())

		wgAdapter.Add(1)
		var sb cache.ScoreBoard
		barCode := progressbar.NewBar(*barIndex+1, "执行["+adapter.Name()+"]", codeCount)
		//updateOneFeature(barAdapter, barCode, adapter, cacheDate, featureDate, op, barIndex)

		mapFeature := treemap.NewWithStringComparator()
		var wg sync.WaitGroup
		//p, _ := ants.NewPoolWithFunc(quotes.POOL_MAX, func(i interface{}) {
		//	v := i.(featureTask)
		//	updateStockFeature(v.wg, v.bar, v.feature, v.securityCode, v.cacheDate, v.featureDate, op, mapFeature)
		//})
		//defer p.Release()

		dataSource := adapter.Factory(featureDate, "")
		parent := coroutine.Context()
		ctx := context.WithValue(parent, cache.KBarIndex, barIndex)
		_ = dataSource.Init(ctx, featureDate)
		for _, code := range allCodes {
			feature := adapter.Factory(cacheDate, code).(factors.Feature)
			if feature.Kind() != factors.FeatureHistory {
				history := smart.GetL5History(code, cacheDate)
				if history != nil {
					feature = feature.FromHistory(*history)
				}
			}
			wg.Add(1)
			go updateStockFeature(&wg, barCode, feature, code, cacheDate, featureDate, op, mapFeature, &sb)
			//go updateStockFeature(&wg, barCode, feature, code, cacheDate, featureDate, op, mapFeature)
			//_ = p.Invoke(featureTask{
			//	wg:           &wg,
			//	bar:          barCode,
			//	feature:      feature,
			//	securityCode: code,
			//	cacheDate:    cacheDate,
			//	featureDate:  featureDate,
			//})
		}
		wg.Wait()
		// 加载缓存
		adapter.Checkout(cacheDate)
		// 合并
		adapter.Merge(mapFeature)
		// 适配器进度条+1
		barAdapter.Add(1)
		wgAdapter.Done()
		fmt.Println(sb.String())
		logger.Infof("%s: %s, end", moduleName, adapter.Name())
	}
	wgAdapter.Wait()
	logger.Infof("%s: all, end", moduleName)
}
