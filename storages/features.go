package storages

import (
	"context"
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/cachel5"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/engine/market"
	"gitee.com/quant1x/engine/smart"
	"gitee.com/quant1x/engine/util"
	"gitee.com/quant1x/gox/coroutine"
	"gitee.com/quant1x/gox/logger"
	"gitee.com/quant1x/gox/progressbar"
	"gitee.com/quant1x/gox/runtime"
	"gitee.com/quant1x/gox/tags"
	"gitee.com/quant1x/gox/text/runewidth"
	"gitee.com/quant1x/gox/util/treemap"
	"gitee.com/quant1x/pkg/tablewriter"
	"os"
	"sync"
	"time"
)

func updateStockFeature(wg *util.RollingWaitGroup, bar *progressbar.Bar, feature factors.Feature, code string, cacheDate, featureDate string, op cache.OpKind, p *treemap.Map, sb *cache.ScoreBoard) {
	defer runtime.CatchPanic()
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
	var metrics []cache.ScoreBoard
	for _, adapter := range adapters {
		logger.Infof("%s: %s, begin", moduleName, adapter.Name())

		wgAdapter.Add(1)
		var sb cache.ScoreBoard
		barCode := progressbar.NewBar(*barIndex+1, "执行["+adapter.Name()+"]", codeCount)
		//updateOneFeature(barAdapter, barCode, adapter, cacheDate, featureDate, op, barIndex)

		mapFeature := treemap.NewWithStringComparator()
		wg := util.NewRollingWaitGroup(5)
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
			sb.Kind = feature.Kind()
			wg.Add(1)
			go updateStockFeature(wg, barCode, feature, code, cacheDate, featureDate, op, mapFeature, &sb)
		}
		wg.Wait()
		// 加载缓存
		adapter.Checkout(cacheDate)
		// 合并
		adapter.Merge(mapFeature)
		// 适配器进度条+1
		barAdapter.Add(1)
		wgAdapter.Done()
		metrics = append(metrics, sb)
		logger.Infof("%s: %s, end", moduleName, adapter.Name())
	}
	wgAdapter.Wait()
	logger.Infof("%s: all, end", moduleName)
	// 输出衡量性能的指标列表
	fmt.Println()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tags.GetHeadersByTags(cache.ScoreBoard{}))
	for _, v := range metrics {
		table.Append(tags.GetValuesByTags(v))
	}
	table.Render()
}
